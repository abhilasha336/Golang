package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"

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

	payload := string(content)

	type MyData struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Hd            string `json:"hd"`
	}

	// Parse the JSON string
	var data MyData
	err = json.Unmarshal([]byte(payload), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// Connection parameters
	host := "localhost"
	port := 5432
	user := "postgres"
	dbname := "test_db"
	sslMode := "disable"

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		host, port, user, dbname, sslMode)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")
	conn, err := pq.NewConnector(connStr)
	if err != nil {
		fmt.Println("connector error")
	}

	idBig, err := strconv.ParseInt(data.ID, 10, 64)
	if err != nil {
		fmt.Println("autherrconv", err)
	}
	email := data.Email
	name := strings.SplitAfter(email, "@")
	exactname := name[0]

	dbc := sql.OpenDB(conn)
	insertQuery := `INSERT INTO users (email, oauth_id, name) VALUES ($1,$2,$3)`

	_, err = dbc.Exec(insertQuery, data.Email, idBig, exactname)
	if err != nil {
		fmt.Println("db insert error", err)
		fmt.Fprintf(w, "User not registered")

	}

	fmt.Fprintf(w, "User registerd successfully")

}
