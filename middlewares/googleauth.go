package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var GOOGLE_ACCOUNT_DOMAIN = []string{"accounts.google.com", "https://accounts.google.com"}

var GOOGLE_TOKEN_VALIDATION_URL = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="

type GoogleAuth struct {
	Domain             []string
	TokenValidationURL string
}

func ReadJsonData(content []uint8) map[string]interface{} {
	var payload map[string]interface{}
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal: ", err)
	}

	return payload
}

func VerifyIss(google_account_domain []string, iss string) bool {
	for _, data := range google_account_domain {
		if iss == data {
			return true
		}
	}

	return false
}

func GoogleTokenValidation(token string) (bool, map[string]interface{}) {
	googleAuth := &GoogleAuth{
		Domain:             GOOGLE_ACCOUNT_DOMAIN,
		TokenValidationURL: GOOGLE_TOKEN_VALIDATION_URL,
	}

	url := fmt.Sprint(googleAuth.TokenValidationURL + token)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	payload := ReadJsonData(data)
	fmt.Println("payload", payload)
	iss := payload["iss"]

	if iss == nil {
		err := make(map[string]interface{})
		err["detail"] = "Invalid Token"
		return false, err
	}

	valid := VerifyIss(googleAuth.Domain, iss.(string))

	return valid, payload
}
