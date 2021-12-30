package http

import (
	"context"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"net/http"
)

type UserCreateRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1"`
	LastName  string `json:"lastName"  validate:"required,min=1"`
	UserId    string `json:"userId"    validate:"required,min=3"`
	Password  string `json:"password"  validate:"required,min=6"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstName,omitempty" validate:"required,min=1"`
	LastName  string `json:"lastName,omitempty"  validate:"required,min=1"`
	UserId    string `json:"userId,omitempty"    validate:"required,min=3"`
	Password  string `json:"password,omitempty"  validate:"required,min=6"`
}

func handleUserOperation(hp HandlerProperties) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO be able to pass context down
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		switch {
		case r.Method == "GET":
			handleGetUser(ctx, hp, w, r)
		case r.Method == "POST":
			handleCreateUser(ctx, hp, w, r)
		case r.Method == "DELETE":
			handleDeleteUser(ctx, hp, w, r)
		case r.Method == "PUT":
			handleUpdateUser(ctx, hp, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleCreateUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	var req UserCreateRequest
	if _, err := parseJSONRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err := hp.validate.Struct(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	userRes, err := hp.userService.CreateUser(ctx, req.FirstName, req.LastName, req.UserId, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respondWithStatus(w, userRes, http.StatusCreated)
}

func handleDeleteUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := hp.validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = hp.userService.DeleteUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respondWithStatus(w, nil, http.StatusNoContent)

}

func handleGetUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := hp.validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user, err := hp.userService.GetUserById(ctx, id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respondWithStatus(w, user, http.StatusOK)
}

func handleUpdateUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	var req UserUpdateRequest
	if _, err := parseJSONRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	id := r.URL.Query().Get("id")
	err := hp.validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = hp.userService.UpdateUser(ctx, id, req.FirstName, req.LastName, req.UserId, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respondWithStatus(w, domain.User{}, http.StatusOK)
}
