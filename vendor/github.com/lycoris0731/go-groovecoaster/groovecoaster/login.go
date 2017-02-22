package groovecoaster

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// tryLogin will attempt login to mypage in GrooveCoaster
// If login was successful, it will return http.Client
func tryLogin() (*http.Client, error) {
	const loginURL = "https://mypage.groovecoaster.jp/sp/login/auth_con.php"

	const nesicaCardID = "NESICA_CARD_ID"
	const nesicaPassword = "NESICA_PASSWORD"

	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	if os.Getenv(nesicaCardID) == "" || os.Getenv(nesicaPassword) == "" {
		return nil, fmt.Errorf("NESiCA ID or password not found")
	}

	v := url.Values{}
	v.Add("nesicaCardId", os.Getenv("NESICA_CARD_ID"))
	v.Add("password", os.Getenv("NESICA_PASSWORD"))

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Jar: jar}

	res, err := client.PostForm(loginURL, v)
	if err != nil {
		return nil, err
	}

	if url := res.Request.URL.String(); strings.Contains(url, "isError=true") {
		return nil, fmt.Errorf("invalid NESiCA Card ID or password")
	} else if strings.Contains(url, "login_stop") {
		return nil, fmt.Errorf("login attempts limit exceed")
	}

	return &client, nil
}
