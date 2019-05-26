package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"os"
	"strings"
	"net/http/httputil"
	"fmt"

	"golang.org/x/oauth2"
	"context"
)

type AuthServer struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURIs []string
	Scopes       []string
}

var state, scope, accessToken string
var port = 9000
var as *AuthServer
var client *Client

func init() {
	as = &AuthServer{
		AuthorizationEndpoint: "http://localhost:9001/authorize",
		TokenEndpoint:         "https://dev-991803.oktapreview.com/oauth2/default/v1/token",
	}

	client = &Client{
		ClientId:     "oauth-client-1",
		ClientSecret: "oauth-client-secret-1",
		RedirectURIs: []string{"http://localhost:9000/callback"},
		Scopes:       []string{"foo", "bar"},
		//Scopes:       []string{"email", "profile", "openid"},
	}

	zerolog.TimeFieldFormat = ""
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = zerolog.New(os.Stdout).With().Caller().Logger()
}

func main() {
	oauthConfig := oauth2.Config{
		ClientID: "0oakuhp8brWUfRhGI0h7",
		ClientSecret: "HNhG1RVIPkqMyZ6PcLR7Ktoxs0geaWoEETRSSy25",
		RedirectURL: "http://localhost:3000/callback",
		Endpoint: oauth2.Endpoint {TokenURL: "https://dev-991803.oktapreview.com/oauth2/default/v1/token"},
	} //ConfigFromJSONの ConfigFromJSONが参考になる
	code := "YCuiQAYpe-fU1KQOotJN"
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Error().Err(err).Msg("hoge")
	} else {
		fmt.Println(token.AccessToken)
		fmt.Println(token.Expiry)
		fmt.Println(token.RefreshToken)
		fmt.Println(token.TokenType)
		fmt.Println(token)
	}

	server := http.Server{Addr: "localhost:9000"}
	http.HandleFunc("/authorize", dumpRequest(handleAuthorize))
	http.HandleFunc("/callback", callback)
	log.Fatal().Err(server.ListenAndServe())
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	endpoint := as.AuthorizationEndpoint

	state = uuid.NewV4().String()

	params := url.Values{
		// rfc6749 required params
		"response_type": {"code"},
		"client_id":     {client.ClientId},
		"redirect_uri":  client.RedirectURIs,

		// rfc6749 optional params
		"scope":         {strings.Join(client.Scopes, " ")},
		"state":         {state},

		// openid params
		//"nonce":         {"this is nonce"},
	}

	if !strings.Contains(as.AuthorizationEndpoint, "?") {
		endpoint = endpoint + "?"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, endpoint+params.Encode(), http.StatusFound)
}

func callback(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(q.Get("code")))
}

type middleware func(next http.HandlerFunc) http.HandlerFunc
func dumpRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump) + "\n")
		next.ServeHTTP(w, r)
	}
}