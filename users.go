package main

import (
	"context"
	"encoding/json"
	"fmt"
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
    
    fmt.Println(user)
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(fmt.Sprintf(`{\n"id": "%s",\n`, user.ID)))
    w.Write([]byte(fmt.Sprintf(`"created_at": "%s",\n`, user.CreatedAt)))
    w.Write([]byte(fmt.Sprintf(`"updated_at": "%s",\n`, user.UpdatedAt)))
    w.Write([]byte(fmt.Sprintf(`"email": "%s",\n}`, user.Email)))
}



