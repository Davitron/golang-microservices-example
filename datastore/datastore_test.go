package main

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"strings"
	"os"
	"github.com/gorilla/handlers"
)

func TestCreate(t *testing.T) {
	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/kv/my-post", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	checkStatus(postResp.Code, http.StatusCreated, t)
}

func TestCreateMissingKey(t *testing.T) {
	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/kv", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)

	// Endpoint is not found
	checkStatus(postResp.Code, http.StatusNotFound, t)
}

func TestGet(t *testing.T) {
	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/kv/my-get", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	checkStatus(postResp.Code, http.StatusCreated, t)

	getReq, err := http.NewRequest("GET", "/kv/my-get", nil)
	checkRequest(err, t)

	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	checkStatus(getResp.Code, http.StatusOK, t)

	if getResp.Body.String() != "my-value" {
		t.Fatal("GET body is not as expeted.")
	}
}

func TestGetMissingKey(t *testing.T) {
	router := handlers.LoggingHandler(os.Stdout, createRouter())

	getReq, err := http.NewRequest("GET", "/kv/my-get-missing", nil)
	checkRequest(err, t)

	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	checkStatus(getResp.Code, http.StatusNotFound, t)
}

func TestDelete(t *testing.T) {
	router := handlers.LoggingHandler(os.Stdout, createRouter())

	// CREATE KEY-VALUE
	postReq, err := http.NewRequest("POST", "/kv/my-delete", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	checkStatus(postResp.Code, http.StatusCreated, t)

	// DELETE KEY-VALUE
	deleteReq, err := http.NewRequest("DELETE", "/kv/my-delete", nil)
	checkRequest(err, t)

	deleteResp := httptest.NewRecorder()
	router.ServeHTTP(deleteResp, deleteReq)
	checkStatus(deleteResp.Code, http.StatusOK, t)

	// GET KEY-VALUE
	getReq, err := http.NewRequest("GET", "/kv/my-delete", nil)
	checkRequest(err, t)

	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	checkStatus(getResp.Code, http.StatusNotFound, t)
}

func checkStatus(actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatal("Server error: Returned ", actual, " instead of ", expected)
	}
}

func checkRequest(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Creating request failed: %+v", err)
	}
}