package service_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/services"
	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(core.HTTPHandleFunc(services.HomeHandler))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var homePageResponse services.APISuccessResponse
	_ = json.NewDecoder(rr.Body).Decode(&homePageResponse)
	expected := services.APISuccessResponse{Detail: "Reetro Application Home Page"}

	if homePageResponse != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAboutPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(core.HTTPHandleFunc(services.AboutHandler))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var aboutPageResponse services.AboutPageResponse
	_ = json.NewDecoder(rr.Body).Decode(&aboutPageResponse)
	expected := services.AboutPageResponse{
		Detail:      "Reetro Application About Page",
		Information: "This application is used for retrospective.",
	}

	if aboutPageResponse != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetRequestUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/about/", nil)
	if err != nil {
		t.Fatal(err)
	}

	user := TestMockUser()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	mockRequestUserStorage := new(MockRequestUserStorage)
	requestUser := services.NewRequestUser(mockRequestUserStorage)
	reqUser := requestUser.GetRequestUser(req)

	assert.Equal(t, reqUser.Id, user.Id)
	assert.Equal(t, reqUser.Email, user.Email)
}
