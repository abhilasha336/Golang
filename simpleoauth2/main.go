package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:3000/google/callback",
	ClientID:     "859339718701-pvvlufnvtjq6b6rorc9h8q1ll3cjo2gp.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

var strCsrf = "abhilash"

func main() {

	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/google/callback", handlerCallback)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("error running server")
		panic(err)

	}
	fmt.Println("server starting")
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html>
			<body>
				<a href="login">Google Sign in</a>
			</body>
		</html>`
	fmt.Fprint(w, html)
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(strCsrf)
	fmt.Println("url is", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handlerCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != strCsrf {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
	token, err := oauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Fprintln(w, err.Error())
		http.Redirect(w, r, "/", http.StatusBadRequest)

	}
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
	fmt.Fprintf(w, string(content))

}
