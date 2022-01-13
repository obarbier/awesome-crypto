package handler

import (
	"encoding/json"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/mongo"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type LoginRequest struct {
	UserId string `json:"userId"`
}

func NewDefaultManageServer() *server.Server {
	// FIXME: this probably should not be there, but instead in a sessionManager class
	session.InitManager(
		session.SetStore(mongo.NewStore("mongodb://127.0.0.1:27017", "users", "session")),
	)
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewFileTokenStore("testDb")) // TODO: implement with REDIS good practice

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("my_secret_key"), jwt.SigningMethodHS512))
	//manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()            // TODO: multi client store// how?
	err := clientStore.Set("222222", &models.Client{ // FIXME
		ID:     "222222",                // FIXME
		Secret: "22222222",              // FIXME
		Domain: "http://localhost:9094", // FIXME
	})
	if err != nil {
		log.Fatalf("failed to create clientStore")
		return nil
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	// TODO: This could be a recovery things
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})
	// TODO: add authorization Handler
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	// TODO: needs to be better
	srv.SetClientInfoHandler(clientInfoHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	return srv
}

func clientInfoHandler(r *http.Request) (clientID, clientSecret string, err error) {
	clientID = "222222"
	clientSecret = "22222222"
	err = nil

	return
}
func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	sessionStore, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := sessionStore.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			err = r.ParseForm()
		}

		sessionStore.Set("ReturnUri", r.Form)
		err = sessionStore.Save()

		w.Header().Set("Location", "/v1/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	sessionStore.Delete("LoggedInUserID")
	err = sessionStore.Save()
	return
}

func handleOauthAuthorize(hp PropertiesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionStore, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//to change the flags on the default logger
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		var form url.Values = make(url.Values)
		if v, ok := sessionStore.Get("ReturnUri"); ok {
			log.Printf("%+v", v)
			getFormData(v, &form)
		}
		if len(form) != 0 {
			r.Form = form
		}

		sessionStore.Delete("ReturnUri")
		err = sessionStore.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = hp.OauthServer.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func handleOauthToken(hp PropertiesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hp.OauthServer.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func handleLogin(_ PropertiesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionStore, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "POST" {
			if r.Form == nil {
				if err := r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			sessionStore.Set("LoggedInUserID", r.Form.Get("username"))
			err = sessionStore.Save()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Location", "/v1/auth")
			w.WriteHeader(http.StatusFound)
			return
		}
		outputHTML(w, r, "/home/obarbier/git/awesome-crypto/user_api/handler/static/login.html") // FIXME

	})
}

func handleTestToken(hp PropertiesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := hp.OauthServer.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		err = e.Encode(data)
		if err != nil {
			return
		}

		portvar := 2021
		log.Printf("Server is running at %d port.\n", portvar)
		log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
		log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")
	})

}

func handleLoginAuth(hp PropertiesHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionStore, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, ok := sessionStore.Get("LoggedInUserID"); !ok {
			w.Header().Set("Location", "/v1/login")
			w.WriteHeader(http.StatusFound)
			return
		}

		outputHTML(w, r, "/home/obarbier/git/awesome-crypto/user_api/handler/static/auth.html") // FIXME
	})
}
func getFormData(m interface{}, formData *url.Values) {
	if formData == nil {
		newForm := make(url.Values)
		formData = &newForm
	}
	for k, v := range m.(map[string]interface{}) {
		switch v.(type) {
		case []interface{}:
			arr := make([]string, len(v.([]interface{})))
			for _, r := range v.([]interface{}) {
				arr = append(arr, strings.TrimSpace(r.(string)))
			}
			formData.Set(k, strings.Join(arr, ""))
		case interface{}:
			formData.Set(k, v.(string))
		default:
			formData.Set(k, v.(string))

		}

	}
}
func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
