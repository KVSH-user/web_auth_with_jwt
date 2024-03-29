Проект, в котором реализована регистрация и авторизация пользоваталей, с генирацией ```JWT``` токена.

Языки, технологии использовались:
1. Golang
2. PostgreSQL
3. REST
4. ROUTER CHI
5. JWT
6. Goose migration

При запуске проекта автоматически накатывается миграция БД, используя ```goose```.

Два ```REST``` маршрута:

- Первый маршрут регистрирует пользователя по данным указанным в параметре запроса.
- Второй маршрут выполняет авторизацию(сверяет логин и пароль, а так же генерирует и выдает JWT access token).

Для регистрации пользователя отправьте ```POST``` запрос на ```/users/signup``` с приведённым телом:

```
{
    "email" : "",
    "password" : ""
}
```

Авторизация

Для авторизации отправьте ```POST``` запрос на ```/users/login``` с приведённым телом:
```
{
    "email" : "",
    "password" : ""
}
```
Пример ответа:
```
{
  "token": <access_token>
}
```
