package handler

import (
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

func TestLoginFlow(t *testing.T) {
	service := TestUserService(t)
	ln, addr := TestServer(t, service)
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			// TODO: handle error
		}
	}(ln)

	token := ""
	resp := testHttpPost(t, token, addr+"/v1/login", map[string]interface{}{
		"userId": "noop",
	})

	testResponseStatus(t, resp, http.StatusOK)
}

func TestOauthAuthorize(t *testing.T) {
	service := TestUserService(t)
	ln, addr := TestServer(t, service)
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			// TODO: handle error
		}
	}(ln)

	token := ""
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	resp := testHttpWithCookiePost(t, token, addr+"/v1/login", map[string]interface{}{
		"userId": "noop",
	}, cookieJar)

	testResponseStatus(t, resp, http.StatusOK)
	resp2 := testHttpFormDataPost(t, token, addr+"/v1/oauth/authorize", url.Values{
		"userId":        {"noop"},
		"redirect_uri":  {"http://localhost:2021/test"},
		"client_id":     {"222222"},
		"response_type": {"token"},
	}, cookieJar)
	testResponseStatus(t, resp2, http.StatusOK)
}
