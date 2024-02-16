package signup

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
	"wishlist_auth/internal/entity"
	apierr "wishlist_auth/internal/http-server/errors"
	resp "wishlist_auth/internal/lib/api/response"
)

type UserAction interface {
	CreateUser(lname, fname, email string, passwordHashed []byte) (int, error)
}

type Response struct {
	resp.Response
	Id int `json:"id,omitempty"`
}

func Create(log *slog.Logger, userAction UserAction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.registration.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req entity.RegistrationRequest

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

		passwordHashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash: ", err)

			return
		}

		id, err := userAction.CreateUser(req.LastName, req.FirstName, req.Email, passwordHashed)
		if err != nil {
			if errors.Is(err, apierr.ErrUsernameTaken) {
				log.Error("username already taken: ", err)
				render.JSON(w, r, resp.Error("username already taken"))
				return
			}

			log.Error("failed to create user: ", err)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("user created")

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
