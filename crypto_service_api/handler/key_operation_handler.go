package handler

import (
	"context"
	"errors"
	"github.com/obarbier/awesome-crypto/common"
	"net/http"
)

// FIXME add validation
type KeyCreate struct {
	KeyName string `json:"key_name"`
	KeyType string `json:"key_type"`
	KeySize int    `json:"key_size"`
}

func keyOperationHandler(hp PropertiesHandler) http.Handler {
	// FIXME since request has a context too.
	ctx := context.Background()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
		case http.MethodPost:
			handleKeyCreate(ctx, hp, w, r)
		}
	})
}

func handleKeyCreate(ctx context.Context, hp PropertiesHandler, w http.ResponseWriter, r *http.Request) {
	var req *KeyCreate = new(KeyCreate)
	if _, err := common.ParseJSONRequest(r, w, req); err != nil {
		common.RespondError(w, http.StatusBadRequest, err)
		return
	}
	err := hp.Validate.Struct(req)
	if err != nil {
		common.RespondError(w, http.StatusBadRequest, err)
		return
	}
	ownerId, ok := r.Context().Value(common.JwtUserDetails).(string)
	if !ok {
		common.RespondError(w, http.StatusBadRequest, errors.New("could not get userId from bearer token"))
		return
	}
	keyResponse, err := hp.UserService.Create(ctx, ownerId, req.KeyName, req.KeyType)
	if err != nil {
		common.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	common.RespondWithStatus(w, keyResponse, http.StatusCreated)
}
