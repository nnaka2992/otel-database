package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"context"
	"log"
	"net/http"
	"io"
	"strconv"

)

var signalChan chan (os.Signal) = make(chan os.Signal, 1)
const url = "http://localhost:8080"

func main() {
	ctx := context.Background()
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:    ":8081",
		Handler: nil,
	}
	go func() {
		http.HandleFunc("/user/", getUserHandler)
		http.HandleFunc("/user/new", postUserAddHandler)
		http.HandleFunc("/user/delete", deleteUserDeleteHandler)
		http.HandleFunc("/user/update", PostUserUpdateHandler)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Receive output from signalChan.
	sig := <-signalChan
	log.Printf("%s signal caught. Graceful Shutdown.", sig)

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}
	log.Print("server exited")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if r.Method != "GET" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "HTTP Method is not valid")
		httpError(w, m)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "Content Type is not valid")
		httpError(w, m)
		return
	}

	userHelper(w, r, url+"/user/", "GET")
}

func postUserAddHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if r.Method != "POST" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "HTTP Method is not valid")
		httpError(w, m)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "Content Type is not valid")
		httpError(w, m)
		return
	}

	userHelper(w, r, url+"/user/new", "POST")
}

func deleteUserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if r.Method != "DELETE" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "HTTP Method is not valid")
		httpError(w, m)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "Content Type is not valid")
		httpError(w, m)
		return
	}

	userHelper(w, r, url+"/user/delete", "DELETE")
}

func PostUserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if r.Method != "POST" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "HTTP Method is not valid")
		httpError(w, m)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		m := fmt.Sprintf("%d Bad Request: %s", http.StatusBadRequest, "Content Type is not valid")
		httpError(w, m)
		return
	}

	userHelper(w, r, url+"/user/update", "POST")
}

func userHelper(w http.ResponseWriter, r *http.Request, url string, method string) {
	// Read request body
	l := r.ContentLength
	body := make([]byte, l)
	_, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func httpError(w http.ResponseWriter, m string) {
	code, err := strconv.Atoi(m[:3])
	if err != nil {
		code = http.StatusInternalServerError
	}
	http.Error(w, m, code)
	log.Printf(m)
}
