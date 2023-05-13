package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type AuthForm struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	GrantType   string `json:"grant_type"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	redirect := "http://127.0.0.1:8080/callback"
	clientID := "?"
	clientSecret := "?"
	basic := clientID + ":" + clientSecret
	encoding := base64.StdEncoding.EncodeToString([]byte(basic))

	listener, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		params := url.Values{}
		params.Add("code", query.Get("code"))
		params.Add("redirect_uri", redirect)
		params.Add("grant_type", "authorization_code")

		request, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(params.Encode()))
		request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encoding))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			panic(err)
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))

		authResponse := &AuthResponse{}
		err = json.Unmarshal(body, &authResponse)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println(authResponse)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		state := "abcdefghijlkmn123"
		scope := "user-read-private user-read-email"
		u, _ := url.Parse("https://accounts.spotify.com/authorize")
		params := url.Values{}
		params.Set("response_type", "code")
		params.Set("client_id", clientID)
		params.Set("scope", scope)
		params.Set("redirect_uri", "http://127.0.0.1:8080/callback")
		params.Set("state", state)
		fmt.Printf("%s?%s", u.String(), params.Encode())
		return
	})
	fmt.Println("listen on: ", listener.Addr().(*net.TCPAddr).Port)
	http.Serve(listener, nil)

}
