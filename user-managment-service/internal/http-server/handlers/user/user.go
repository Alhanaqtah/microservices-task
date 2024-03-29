package user

import (
	"errors"
	"log/slog"
	"net/http"
	"user-managment-service/internal/config"
	"user-managment-service/internal/lib/jwt"
	"user-managment-service/internal/lib/logger/sl"
	resp "user-managment-service/internal/lib/response"
	"user-managment-service/internal/models"
	service "user-managment-service/internal/service/user"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	gojwt "github.com/golang-jwt/jwt/v5"
)

type Service interface {
	UserByUUID(uuid string) (*models.User, error)
	PatchUser(uuid string, user *models.User) (*models.User, error)
	Delete(uuid string) error
}

type Handler struct {
	log      *slog.Logger
	service  Service
	tokenCfg config.Token
}

func New(log *slog.Logger, service Service, tokenCfg config.Token) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		tokenCfg: tokenCfg,
	}
}

func (h *Handler) Register() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/me", h.get)
		r.Patch("/me", h.patch)
		r.Delete("/me", h.delete)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.get"

	log := h.log.With(slog.String("op", op))

	// Retrive user id
	tokenString := jwtauth.TokenFromHeader(r)
	token, err := gojwt.Parse(tokenString, func(t *gojwt.Token) (interface{}, error) {
		return []byte(h.tokenCfg.JWT.Secret), nil
	})
	if err != nil {
		log.Error("failed to get user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	claims := token.Claims.(gojwt.MapClaims)

	uuid, err := jwt.GetClaim(claims, "sub")
	if err != nil {
		log.Error("failed to get user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	user, err := h.service.UserByUUID(uuid)
	if err != nil {
		log.Error("failed to get user", sl.Error(err))
		if errors.As(err, &service.ErrUserNotFound) {
			render.JSON(w, r, resp.Err("user not found"))
			return
		}
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, user)
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.patch"

	log := h.log.With(slog.String("op", op))

	// Retrive user id
	tokenString := jwtauth.TokenFromHeader(r)
	token, err := gojwt.Parse(tokenString, func(t *gojwt.Token) (interface{}, error) {
		return []byte(h.tokenCfg.JWT.Secret), nil
	})
	if err != nil {
		log.Error("failed to patch user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	claims := token.Claims.(gojwt.MapClaims)

	uuid, err := jwt.GetClaim(claims, "sub")
	if err != nil {
		log.Error("failed to patch user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	var user models.User
	err = render.DecodeJSON(r.Body, &user)
	if err != nil {
		log.Error("failed to patch user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	u, err := h.service.PatchUser(uuid, &user)
	if err != nil {
		log.Error("failed to patch user", sl.Error(err))
		if errors.As(err, &service.ErrNoFieldsToUpdate) {
			render.JSON(w, r, resp.Err("no fields to update"))
			return
		}
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, u)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.delete"

	log := h.log.With(slog.String("op", op))

	// Retrive user id
	tokenString := jwtauth.TokenFromHeader(r)
	token, err := gojwt.Parse(tokenString, func(t *gojwt.Token) (interface{}, error) {
		return []byte(h.tokenCfg.JWT.Secret), nil
	})
	if err != nil {
		log.Error("failed to delete user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	claims := token.Claims.(gojwt.MapClaims)

	uuid, err := jwt.GetClaim(claims, "sub")
	if err != nil {
		log.Error("failed to delete user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	err = h.service.Delete(uuid)
	if err != nil {
		log.Error("failed to delete atch user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, resp.Ok())
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) patchByID(w http.ResponseWriter, r *http.Request) {

}
