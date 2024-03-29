package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kokweikhong/go-playground/auth-nextjs-golang/backend/internal/middleware"
)

type UserController interface {
	ListUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type userController struct{}

func NewUserController() UserController {
	return &userController{}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var users = []*User{
	{1, "John Doe", "johndoe@gmail.com", "123456", "admin"},
	{2, "Jane Doe", "janedoe@gmail.com", "123456", "user"},
}

func (c *userController) ListUsers(w http.ResponseWriter, r *http.Request) {
	slog.Info("ListUsers")
	data := map[string]interface{}{
		"data": users,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *userController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	users = append(users, &user)
	json.NewEncoder(w).Encode(user)
}

func (c *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	for _, user := range users {
		if userID == fmt.Sprintf("%d", user.ID) {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (c *userController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user *User
	userID := chi.URLParam(r, "userId")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, u := range users {
		if userID == fmt.Sprintf("%d", u.ID) {
			u.Username = user.Username
			u.Email = user.Email
			u.Password = user.Password
			json.NewEncoder(w).Encode(u)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)

}

func (c *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	for i, user := range users {
		if userID == fmt.Sprintf("%d", user.ID) {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	for _, u := range users {
		if strings.EqualFold(user.Email, u.Email) && strings.EqualFold(user.Password, u.Password) {
			token, payload, err := middleware.NewJWTClaims(user.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			rToken, rPayload, err := middleware.NewRefreshToken(user.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var data = struct {
				AccessToken       string `json:"accessToken"`
				AccessTokenExpiry int64  `json:"accessTokenExpiry"`
				RefreshToken      string `json:"refreshToken"`
				RereshTokenExpiry int64  `json:"refreshTokenExpiry"`
				*User
			}{
				AccessToken:       token,
				AccessTokenExpiry: payload.Expiry,
				RefreshToken:      rToken,
				RereshTokenExpiry: rPayload.Expiry,
				User:              u,
			}

			json.NewEncoder(w).Encode(data)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)

}

func (c *userController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from body
	type RefreshToken struct {
		Username string `json:"username"`
		Token    string `json:"refreshToken"`
	}

	var rt RefreshToken

	if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate refresh token
	valid, err := middleware.ValidateJWTToken(rt.Token)
	if err != nil || !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if valid, generate new access token
	tokenString, payload, err := middleware.NewJWTClaims(rt.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return new access token
	var data = struct {
		AccessToken       string `json:"accessToken"`
		AccessTokenExpiry int64  `json:"accessTokenExpiry"`
	}{
		AccessToken:       tokenString,
		AccessTokenExpiry: payload.Expiry,
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
