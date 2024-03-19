package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"user-managment-service/internal/config"
	"user-managment-service/internal/lib/logger/sl"
	resp "user-managment-service/internal/lib/response"
	"user-managment-service/internal/models"
	service "user-managment-service/internal/service/auth"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Service interface {
	SignUp(username, email, password string) (string, error)
	Login(username, password string) (string, string, error)
	RefreshTokens(token string) (string, string, error)
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
		r.Post("/signup", h.signup)
		r.Post("/login", h.login)
		r.Post("/refresh-token", h.refreshToken)
		r.Post("/reset-password", h.resetPassword)
	}
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.signup"

	log := h.log.With(slog.String("op", op))

	var user models.User
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		log.Error("failed to signup user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		log.Debug("failed to signup user: invalid credentials")
		render.JSON(w, r, resp.Err("invalid credentials"))
		return
	}

	uuid, err := h.service.SignUp(user.Username, user.Email, user.Password)
	if err != nil {
		log.Debug("failed to signup user", sl.Error(err))
		if errors.As(err, &service.ErrUserExists) {
			render.JSON(w, r, resp.Err("user already exists"))
			return
		}
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, models.User{UUID: uuid})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.login"

	log := h.log.With(slog.String("op", op))

	var user models.User
	render.DecodeJSON(r.Body, &user)

	if user.Username == "" && user.Password == "" {
		log.Debug("invalid credentials")
		render.JSON(w, r, resp.Err("invalid credentials"))
	}

	accessToken, refreshToken, err := h.service.Login(user.Username, user.Password)
	if err != nil {
		log.Error("failed to login user", sl.Error(err))
		if errors.As(err, &service.ErrUserNotFound) {
			render.JSON(w, r, resp.Err("user not found"))
			return
		}
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, resp.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.refreshToken"

	log := h.log.With(slog.String("op", op))

	type RefreshTokenRequest struct {
		Token string `json:"refreshToken"`
	}

	var oldRefreshToken RefreshTokenRequest
	err := render.DecodeJSON(r.Body, &oldRefreshToken)
	if err != nil {
		log.Debug("failed to refresh tokens", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	accessToken, refreshToken, err := h.service.RefreshTokens(oldRefreshToken.Token)
	if err != nil {
		log.Error("failed to refresh tokens", sl.Error(err))
		if errors.As(err, &service.ErrTokenRevoked) {
			render.JSON(w, r, resp.Err("token revoked"))
			return
		}
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	render.JSON(w, r, resp.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {

}
