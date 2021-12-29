package http

import (
	"net"
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
		"password":  "foo",
	})
	testResponseStatus(t, resp, 204)

}
