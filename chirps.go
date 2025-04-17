package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/MalikL2005/http_server/internal/database"
	"github.com/google/uuid"
)


var unallowed_words []string = []string{
    "kerfuffle",
    "sharbert",
    "fornax",
}

type Chirp struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    UserID    uuid.UUID `json:"user_id"`
    Body      string `json:"body"`
}


func handleValidation (writer http.ResponseWriter, res *http.Request){
    defer res.Body.Close()
    type validation struct {
        Body string `json:"body"`
    }
    body, err := io.ReadAll(res.Body)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{"error":"Could not read response}"`))
        return
    }
    val := validation{}
    err = json.Unmarshal(body, &val)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{"error":"Could not process json-data"}`))
        return
    }
    if len(val.Body) > 140 {
        writer.WriteHeader(http.StatusBadRequest)
        writer.Write([]byte(`{"error": "Chirp is too long"}`))
        return
    }
    for _, word := range(unallowed_words){
        if index := strings.Index(strings.ToLower(val.Body), word); index >= 0{
            val.Body = val.Body[:index] + "****" + val.Body[index+len(word):]
        }
    }

    writer.WriteHeader(http.StatusOK)
    writer.Write(fmt.Appendf([]byte{}, `{"cleaned_body": "%s"}`, val.Body))
}




func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request){
    type json_chirp struct {
        Body    string `json:"body"`
        UserID  uuid.UUID `json:"user_id"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write(fmt.Appendf([]byte{}, `"error": "%s"`, err))
        return
    }

    json_data := map[string]string{}
    err = json.Unmarshal(body, &json_data)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write(fmt.Appendf([]byte{}, `"error": "%s"`, err))
        return
    }

    user_id, ok := json_data["user_id"]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`"error": "no user id found"`))
        return
    }

    parsed_uid, err := uuid.Parse(user_id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`"error": "invalid user id"`))
        return
    }

    chirpBody, ok := json_data["body"]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`"error": "no body found"`))
        return
    }

    chirp, err := cfg.DB.CreateChirp(context.Background(), database.CreateChirpParams{UserID: parsed_uid, Body: chirpBody})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`"error": "could not create DB entry"`))
        return
    }

    w.WriteHeader(http.StatusCreated)
    json_res, err := json.Marshal(&chirp)
    if err != nil {
        w.Write([]byte(`{"created": "sucessfull", "error": "Marshal failed"}`))
        return
    }
    
    w.Write([]byte(json_res))
}




