package http

import (
	"context"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"net/http"
)

type UserCreateRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserId    string `json:"userId"`
	Password  string `json:"password"`
}

type UserCreateResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserId    string `json:"userId"`
	Password  string `json:"password"`
}

func handleUserOperation(service domain.IUserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO be able to pass context down
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		switch {
		case r.Method == "GET":
			handleGetUser(ctx, service, w, r)
		case r.Method == "POST":
			handleCreateUser(ctx, service, w, r)
		case r.Method == "DELETE":
			handleDeleteUser(ctx, service, w, r)
		case r.Method == "PUT":
			handleUpdateUser(ctx, service, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleCreateUser(ctx context.Context, service domain.IUserService, w http.ResponseWriter, r *http.Request) {
	var req UserCreateRequest
	if _, err := parseJSONRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// TODO: in a hexagonal architecture should input validation live in service handler
	// TODO: Validation
	firstName := req.FirstName
	if len(firstName) == 0 {
		// TODO : error
	}

	lastName := req.LastName
	if len(lastName) == 0 {
		// TODO : error
	}
	userId := req.UserId
	if len(userId) == 0 {
		// TODO : error
	}
	password := req.Password
	if len(password) < 6 {
		// TODO : error
	}

	_, err := service.CreateUser(ctx, firstName, lastName, userId, password)
	if err != nil {
		// TODO: Handle errors
		return
	}
}

func handleDeleteUser(ctx context.Context, service domain.IUserService, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}

func handleGetUser(ctx context.Context, service domain.IUserService, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}

func handleUpdateUser(ctx context.Context, service domain.IUserService, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement me
	panic("Implement me")
}
