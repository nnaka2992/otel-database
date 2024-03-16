package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func httpError(w http.ResponseWriter, m string) {
	code, err := strconv.Atoi(m[:3])
	if err != nil {
		code = http.StatusInternalServerError
	}
	http.Error(w, m, code)
	log.Printf(m)
}

func readJson(r *http.Request) (map[string]interface{}, error) {
	l := r.ContentLength
	body := make([]byte, l)
	_, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
	}
	var params map[string]interface{}
	err = json.Unmarshal(body[:l], &params)
	if err != nil {
		return nil, fmt.Errorf("%d Bad Request: %s", http.StatusBadRequest, err)
	}
	return params, nil
}
