package auth

import (
	"log/slog"
	"net/http"
	"user-managment-service/internal/config"
	"user-managment-service/internal/lib/logger/sl"
	resp "user-managment-service/internal/lib/response"
	"user-managment-service/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	FindUser(username string, passHash []byte) (bool, error)
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
	err := render.DecodeJSON(r.Body, user)
	if err != nil {
		log.Error("failed to signup user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	if user.Username != "" && user.Password != "" {
		log.Debug("failed to signup user: invalid credentials")
		render.JSON(w, r, resp.Err("invalid credentials"))
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to signup user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}

	found, err := h.service.FindUser(user.Username, passHash)
	if err != nil {
		log.Error("failed to signup user", sl.Error(err))
		render.JSON(w, r, resp.Err("internal error"))
		return
	}
	if !found {
		log.Debug("failed to signup user: user not found")
		render.JSON(w, r, resp.Err("user not found"))
		return
	}

}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {

}
