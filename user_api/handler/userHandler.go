package handler

import (
	"context"
	"net/http"
)

type UserCreateRequest struct {
	FirstName string `json:"firstName" validate:"required,not-blank"`
	LastName  string `json:"lastName"  validate:"required,not-blank"`
	UserId    string `json:"userId"    validate:"required,not-blank"`
	Password  string `json:"password"  validate:"required,not-blank"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstName,omitempty" validate:"required,not-blank"`
	LastName  string `json:"lastName,omitempty"  validate:"required,not-blank"`
	UserId    string `json:"userId,omitempty"    validate:"required,not-blank"`
	Password  string `json:"password,omitempty"  validate:"required,not-blank"`
}

func handleUserOperation(hp PropertiesHandler) http.Handler {
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

func handleCreateUser(ctx context.Context, hp PropertiesHandler, w http.ResponseWriter, r *http.Request) {
	var req *UserCreateRequest = new(UserCreateRequest)
	if _, err := parseJSONRequest(r, w, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	err := hp.Validate.Struct(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	userRes, err := hp.UserService.CreateUser(ctx, req.FirstName, req.LastName, req.UserId, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithStatus(w, userRes, http.StatusCreated)
}

func handleDeleteUser(ctx context.Context, hp PropertiesHandler, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := hp.Validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = hp.UserService.DeleteUser(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithStatus(w, nil, http.StatusNoContent)

}

func handleGetUser(ctx context.Context, hp PropertiesHandler, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := hp.Validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user, err := hp.UserService.GetUserById(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithStatus(w, user, http.StatusOK)
}

func handleUpdateUser(ctx context.Context, hp PropertiesHandler, w http.ResponseWriter, r *http.Request) {
	var req UserUpdateRequest
	if _, err := parseJSONRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	id := r.URL.Query().Get("id")
	err := hp.Validate.Var(id, "required,uuid")
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = hp.UserService.UpdateUser(ctx, id, req.FirstName, req.LastName, req.UserId, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithStatus(w, nil, http.StatusOK)
}
