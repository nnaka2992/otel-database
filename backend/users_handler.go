package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/nnaka2992/otel-database/backend/gen/sqlc"
)

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

	// Read request body
	params, err := readJson(r)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	if params["Name"] == nil || params["Email"] == nil || params["Age"] == nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, "Invalid Input")
		httpError(w, m)
		return
	}

	ctx := context.Background()
	age, err := strconv.Atoi(params["Age"].(string))
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, "Invalid Input")
		httpError(w, m)
		return
	}
	u, err := query.CreateUser(ctx, db.CreateUserParams{
		Name:  params["Name"].(string),
		Email: params["Email"].(string),
		Age:   int32(age),
	})
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("User Created: %v\n", u)))
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

	// Read request body
	params, err := readJson(r)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	ctx := context.Background()
	var u db.User
	if params["ID"] != nil && params["ID"] != 0 {
		u, err = query.DeleteUserByID(ctx, params["ID"].(int32))
	} else if params["Email"] != nil && params["Email"] != "" {
		u, err = query.DeleteUserByEmail(ctx, params["Email"].(string))
	} else {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, "Invalid Input")
		httpError(w, m)
		return
	}
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}

	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("User Deleted: %v\n", u)))
}

func postUserUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	// Read request body
	params, err := readJson(r)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	ctx := context.Background()
	var u db.User
	if params["ID"] != nil && params["ID"] != 0 {
		u, err = query.GetUserByID(ctx, params["ID"].(int32))
	} else if params["Email"] != nil && params["Email"] != "" {
		u, err = query.GetUserByEmail(ctx, params["Email"].(string))
	} else {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, "User not found")
		httpError(w, m)
		return
	}
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}

	uu := db.UpdateUserParams{
		ID: u.ID,
	}
	if params["Name"] != nil && params["Name"] != "" {
		uu.Name = u.Name
	} else {
		uu.Name = params["Name"].(string)
	}
	if params["Email"] != nil && params["Email"] == "" {
		uu.Email = u.Email
	} else {
		uu.Email = params["Email"].(string)
	}
	if params["Age"] != nil && params["Age"] == "" {
		uu.Age = u.Age
	} else {
		age, err := strconv.Atoi(params["Age"].(string))
		if err != nil {
			m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
			httpError(w, m)
			return
		}
		uu.Age = int32(age)
	}
	u, err = query.UpdateUser(ctx, uu)
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("User Updated: %v\n", u)))
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

	// Read request body
	params, err := readJson(r)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	ctx := context.Background()
	var u db.User
	if params["ID"] != nil && params["ID"] != 0 {
		u, err = query.GetUserByID(ctx, params["ID"].(int32))
	} else if params["Email"] != nil && params["Email"] != "" {
		u, err = query.GetUserByEmail(ctx, params["Email"].(string))
	} else {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, "User not found")
		httpError(w, m)
		return
	}
	if err != nil {
		m := fmt.Sprintf("%d Internal Server Error: %s", http.StatusInternalServerError, err)
		httpError(w, m)
		return
	}
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("User Updated: %v\n", u)))
}
