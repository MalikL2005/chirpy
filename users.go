package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)


func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request){
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusTeapot)
        w.Write([]byte(`"error": "You are a teapot and your json-data could not be serialized"`))
        return
    }
    var jsonData map[string]string
    err = json.Unmarshal(body, &jsonData)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`"error": "could not unmarshal json"`))
        return
    }

    email, ok := jsonData["email"]
    if !ok {
        w.WriteHeader(http.StatusTeapot)
        w.Write([]byte(`"error": "You are a teapot and your request does not contain an email"`))
        return
    }

    user, err := cfg.DB.CreateUser(context.Background(), email)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`"error": "Could not create user in DB"`))
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json_user, err := json.Marshal(&user)
    if err != nil {
        w.Write([]byte(`{"created": "sucessfull", "error": "Marshal failed"}`))
        return
    }
    w.Write(json_user)
}



