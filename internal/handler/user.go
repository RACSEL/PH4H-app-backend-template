package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ips-lacpass-backend/internal/core"
	errors2 "ips-lacpass-backend/internal/errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserService core.UserService
}

func NewUserHandler(s core.UserService) *UserHandler {
	return &UserHandler{UserService: s}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func translateError(body []map[string]interface{}) []map[string]interface{} {
	for _, m := range body {
		if val, ok := m["error"]; ok {
			m["error"] = toSnakeCase(val.(string))
		}
	}
	return body
}

type UserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Locale       string `json:"locale"`
	DocumentType string `json:"document_type"`
	Identifier   string `json:"identifier"`
}

type userCreationRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	FirstName       string `json:"first_name,omitempty" validate:"required"`
	LastName        string `json:"last_name,omitempty" validate:"required"`
	Locale          string `json:"locale" validate:"required,oneof=es en pt-br"`
	DocumentType    string `json:"document_type" validate:"required,oneof=passport identifier"`
	Identifier      string `json:"identifier" validate:"required"`
}

type userUpdateRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" validate:"required"`
}

// Create User godoc
//
//	@Summary		Register a new Keycloak user
//	@Description	Register a new Keycloak user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//
//	@Param			user	body		core.UserRequest	true	"New user parameters"
//
//	@Success		201		{object}	UserResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/users [post]
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body userCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		var verr []map[string]string
		for _, err := range err.(validator.ValidationErrors) {
			name := toSnakeCase(err.Field())
			switch err.Tag() {
			case "required":
				verr = append(verr, map[string]string{
					"error":             fmt.Sprintf("missing_%s", name),
					"error_description": fmt.Sprintf("Missing required field: %s", err.Field()),
				})
			case "email":
				verr = append(verr, map[string]string{
					"error":             fmt.Sprintf("invalid_%s", name),
					"error_description": fmt.Sprintf("Invalid %s type", strings.ReplaceAll(name, "_", " ")),
				})
			case "oneof":
				verr = append(verr, map[string]string{
					"error":             fmt.Sprintf("invalid_%s", name),
					"error_description": fmt.Sprintf("Invalid %s. Must be either %s", strings.ReplaceAll(name, "_", " "), strings.ReplaceAll(err.Param(), " ", " or ")),
				})
			}
		}
		res, err := json.Marshal(verr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if body.Password != body.PasswordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		res, err := json.Marshal([]map[string]string{
			{
				"error":             "invalid_password_confirm",
				"error_description": "Password and password confirmation do not match",
			},
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}

	user, err := u.UserService.CreateUser(r.Context(), core.UserRequest{
		Email:           body.Email,
		Password:        body.Password,
		PasswordConfirm: body.PasswordConfirm,
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		Locale:          body.Locale,
		DocumentType:    core.AllowedDocumenTypes[body.DocumentType],
		Identifier:      body.Identifier,
	})

	if err != nil {
		var cuErr *errors2.HttpError
		if errors.As(err, &cuErr) {
			res, err := json.Marshal(translateError(cuErr.Body))
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if cuErr.StatusCode == 409 {
				// TODO Conflict error, user already exists. Cannot give this details to the user.
				res = []byte(`[{"error":"user_already_exists","error_description":"User already exists"}]`)
			}
			w.WriteHeader(cuErr.StatusCode)
			w.Write(res)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(
		&UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Locale:       user.Locale,
			DocumentType: string(user.DocumentType),
			Identifier:   user.Identifier,
		})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// UpdateUser godoc
//
//	@Summary		Update user profile
//	@Description    Update user profile. Only firs name, last name for now
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//
//	@Param			user	body		core.UserUpdateRequest	true	"New user details"
//
// @Security		ApiKeyAuth
//
//	@Success		200		{object}	UserResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/users/auth/update [put]
func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body userUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		var verr []map[string]string
		for _, err := range err.(validator.ValidationErrors) {
			name := toSnakeCase(err.Field())
			switch err.Tag() {
			case "required":
				verr = append(verr, map[string]string{
					"error":             fmt.Sprintf("missing_%s", name),
					"error_description": fmt.Sprintf("Missing required field: %s", err.Field()),
				})
			}
		}
		res, err := json.Marshal(verr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	user, err := u.UserService.UpdateUser(r.Context(), core.UserUpdateRequest{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	})

	if err != nil {
		var cuErr *errors2.HttpError
		if errors.As(err, &cuErr) {
			res, err := json.Marshal(translateError(cuErr.Body))
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(cuErr.StatusCode)
			w.Write(res)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(
		&UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Locale:       user.Locale,
			DocumentType: string(user.DocumentType),
			Identifier:   user.Identifier,
		})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
