package signin

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
	"time"
	"wishlist_auth/internal/entity"
	apierr "wishlist_auth/internal/http-server/errors"
	resp "wishlist_auth/internal/lib/api/response"
	"wishlist_auth/internal/lib/jwttoken"
)

const secret = "dhaw7dyaw8"

type UserLogin interface {
	Login(email string) ([]byte, int, error)
}

type Response struct {
	Token string `json:"token"`
}

func SignIn(log *slog.Logger, userLogin UserLogin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.SignIn"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req entity.AuthRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		passwordHashed, id, err := userLogin.Login(req.Email)
		if err != nil {
			if errors.Is(err, apierr.ErrIncorrectEmail) {
				log.Error("incorrect email: ", err)
				render.JSON(w, r, resp.Error("incorrect credentials"))
				return
			}
			log.Error("failed to get password: ", err)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		if err := bcrypt.CompareHashAndPassword(passwordHashed, []byte(req.Password)); err != nil {
			log.Error("invalid password: ", err)
			render.JSON(w, r, resp.Error("invalid credential"))
			return
		}

		token, err := jwttoken.Generate(secret, id, time.Hour*24)

		log.Info("user successful login")

		render.JSON(w, r, Response{
			Token: token,
		})

	}
}
