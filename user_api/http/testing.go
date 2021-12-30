package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/obarbier/awesome-crypto/user_api/adapters"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestServer(tb testing.TB, service domain.IUserService) (net.Listener, string) {
	ln, addr := TestListener(tb)
	TestServerWithListener(tb, ln, addr, service)
	return ln, addr
}

func TestListener(tb testing.TB) (net.Listener, string) {
	fail := func(format string, args ...interface{}) {
		panic(fmt.Sprintf(format, args...))
	}
	if tb != nil {
		fail = tb.Fatalf
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fail("err: %s", err)
	}
	addr := "http://" + ln.Addr().String()
	return ln, addr
}

func TestServerWithListener(tb testing.TB, ln net.Listener, addr string, service domain.IUserService) {
	//ip, _, _ := net.SplitHostPort(ln.Addr().String())

	// Create a muxer to handle our requests so that we can authenticate
	// for tests.
	props := HandlerProperties{
		userService:  service,
		recoveryMode: false,
		validate:     validator.New(),
	}
	TestServerWithListenerAndProperties(tb, ln, addr, service, props)
}

func TestServerWithListenerAndProperties(tb testing.TB, ln net.Listener, addr string, service domain.IUserService, props HandlerProperties) {
	// Create a muxer to handle our requests so that we can authenticate
	// for tests.
	mux := http.NewServeMux()
	//mux.Handle("/_test/auth", http.HandlerFunc(testHandleAuth))
	mux.Handle("/", Handler(props))

	server := &http.Server{
		Addr:     ln.Addr().String(),
		Handler:  mux,
		ErrorLog: log.Default(), // TODO: configure logger
	}
	go func() {
		err := server.Serve(ln)
		if err != nil {
			// TODO : Handle error
		}
	}()
}

// TestUserService returns a pure in-memory, uninitialized core for testing.
func TestUserService(t *testing.T) domain.IUserService {
	repo, err := adapters.NewMongoRepository() // FIXME: add in memory dependency
	if err != nil {
		t.Fatal(err)
	}
	return domain.NewUserService(repo)
}

func testResponseStatus(t *testing.T, resp *http.Response, code int) {
	t.Helper()
	if resp.StatusCode != code {
		body := new(bytes.Buffer)
		_, err1 := io.Copy(body, resp.Body)
		if err1 != nil {
			return
		}
		err2 := resp.Body.Close()
		if err2 != nil {
			return
		}

		t.Fatalf(
			"Expected status %d, got %d. Body:\n\n%s",
			code, resp.StatusCode, body.String())
	}
}

func testHttpGet(t *testing.T, token string, addr string) *http.Response {
	loggedToken := token
	if len(token) == 0 {
		loggedToken = "<empty>"
	}
	t.Logf("Token is %s", loggedToken)
	return testHttpData(t, "GET", token, addr, nil, false, 0)
}

func testHttpDelete(t *testing.T, token string, addr string) *http.Response {
	return testHttpData(t, "DELETE", token, addr, nil, false, 0)
}

// Go 1.8+ clients redirect automatically which breaks our 307 standby testing
func testHttpDeleteDisableRedirect(t *testing.T, token string, addr string) *http.Response {
	return testHttpData(t, "DELETE", token, addr, nil, true, 0)
}

func testHttpPostWrapped(t *testing.T, token string, addr string, body interface{}, wrapTTL time.Duration) *http.Response {
	return testHttpData(t, "POST", token, addr, body, false, wrapTTL)
}

func testHttpPost(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "POST", token, addr, body, false, 0)
}

func testHttpPut(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "PUT", token, addr, body, false, 0)
}

// Go 1.8+ clients redirect automatically which breaks our 307 standby testing
func testHttpPutDisableRedirect(t *testing.T, token string, addr string, body interface{}) *http.Response {
	return testHttpData(t, "PUT", token, addr, body, true, 0)
}

func testHttpData(t *testing.T, method string, token string, addr string, body interface{}, disableRedirect bool, wrapTTL time.Duration) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		if err := enc.Encode(body); err != nil {
			t.Fatalf("err:%s", err)
		}
	}

	req, err := http.NewRequest(method, addr, bodyReader)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get the address of the local listener in order to attach it to an Origin header.
	// This will allow for the testing of requests that require CORS, without using a browser.
	hostURLRegexp, _ := regexp.Compile("http[s]?://.+:[0-9]+")
	req.Header.Set("Origin", hostURLRegexp.FindString(addr))

	req.Header.Set("Content-Type", "application/json")

	if wrapTTL > 0 {
		req.Header.Set("X-Vault-Wrap-TTL", wrapTTL.String())
	}

	//if len(token) != 0 {
	//	req.Header.Set(consts.AuthHeaderName, token)
	//}

	client := cleanhttp.DefaultClient()
	client.Timeout = 60 * time.Second

	// From https://github.com/michiwend/gomusicbrainz/pull/4/files
	defaultRedirectLimit := 30

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if disableRedirect {
			return fmt.Errorf("checkRedirect disabled for test")
		}
		if len(via) > defaultRedirectLimit {
			return fmt.Errorf("%d consecutive requests(redirects)", len(via))
		}
		if len(via) == 0 {
			// No redirects
			return nil
		}
		// TODO: handle
		//// mutate the subsequent redirect requests with the first Header
		//if token := via[0].Header.Get(consts.AuthHeaderName); len(token) != 0 {
		//	req.Header.Set(consts.AuthHeaderName, token)
		//}
		return nil
	}

	resp, err := client.Do(req)
	if err != nil && !strings.Contains(err.Error(), "checkRedirect disabled for test") {
		t.Fatalf("err: %s", err)
	}

	return resp
}
