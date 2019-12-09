package main

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	json := `
	{"url":"http://omotenashi-customer-site.com/","key":"uboC0Et70nU5ShPlAwOnmMZS5xt1","once":[]}
`
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(json)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, r)
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
}
