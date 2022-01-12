package handler

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"
)

func TestUserCreate(t *testing.T) {
	service := TestUserService(t)
	ln, addr := TestServer(t, service)
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			// TODO: handle error
		}
	}(ln)
	token := ""
	resp := testHttpPost(t, token, addr+"/v1/user", map[string]interface{}{
		"firstName": "noop",
		"lastName":  "foo",
		"userId":    "foo",
		"password":  "password123",
	})

	testResponseStatus(t, resp, http.StatusCreated)

	entity := getResponse(resp)
	id := entity["id"].(string)

	//update
	resp2 := testHttpPut(t, token, addr+"/v1/user?id="+id, map[string]interface{}{
		"userId":   "updatedId",
		"password": "passwordUpdated",
	})
	testResponseStatus(t, resp2, http.StatusOK)
	// get
	resp3 := testHttpGet(t, token, addr+"/v1/user?id="+id)
	testResponseStatus(t, resp3, http.StatusOK)
	// delete
	resp4 := testHttpDelete(t, token, addr+"/v1/user?id="+id)
	testResponseStatus(t, resp4, http.StatusNoContent)

}

func TestValidation(t *testing.T) {
	service := TestUserService(t)
	ln, addr := TestServer(t, service)
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			// TODO: handle error
		}
	}(ln)
	token := ""
	resp := testHttpPost(t, token, addr+"/v1/user", map[string]interface{}{
		"firstName": "noop",
		"userId":    "foo",
	})

	testResponseStatus(t, resp, http.StatusBadRequest)
}

func getResponse(resp *http.Response) map[string]interface{} {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bodyByte, &jsonData)
	if err != nil {
		return nil
	}

	return jsonData
}
