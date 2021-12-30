package http

import (
	"context"
	"net/http"
)

type UserCreateRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
	UserId    string `json:"userId"    validate:"required"`
	Password  string `json:"password"  validate:"required"`
}

type UserCreateResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserId    string `json:"userId"`
	Password  string `json:"password"`
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

	_, err = hp.userService.CreateUser(ctx, req.FirstName, req.LastName, req.UserId, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
}

func handleDeleteUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}

func handleGetUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}

func handleUpdateUser(ctx context.Context, hp HandlerProperties, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}
