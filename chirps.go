package main

import (
	"context"
	"encoding/json"
	"errors"
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


func validateChirp (body string) (string, error){
    type validation struct {
        Body string `json:"body"`
    }
    val := validation{}
    err := json.Unmarshal([]byte(body), &val)
    if err != nil {
        return "", err
    }
    if len(val.Body) > 140 {
        return "", errors.New(fmt.Sprintf("Max length is 140 chars. Chirp is %d characters long", len(body)))
    }
    for _, word := range(unallowed_words){
        if index := strings.Index(strings.ToLower(val.Body), word); index >= 0{
            val.Body = val.Body[:index] + "****" + val.Body[index+len(word):]
        }
    }
    return val.Body, nil
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
    fmt.Println(json_data["body"])

    json_data["body"], err = validateChirp(string(body))
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

    fmt.Println(json_data)
    fmt.Println(user_id)

    parsed_uid, err := uuid.Parse(string(user_id))
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
    
    chirpFormatted := Chirp{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        UserID: chirp.ID,
        Body: chirp.Body,
    }

    jsonResponse, err := json.Marshal(chirpFormatted)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`"error": "Marshal json response failed"`))
    }

    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}



func (cfg *apiConfig) handleGetAllChirps (w http.ResponseWriter, r *http.Request){
    chirps, err := cfg.DB.GetAllChirps(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(fmt.Appendf([]byte{}, `{"error": %s}`, err))
    }

    chirpsFormatted := make([]Chirp, 0)
    for _, chirp := range chirps{
        chirpsFormatted = append(chirpsFormatted, Chirp{
            ID: chirp.ID, 
            CreatedAt: chirp.CreatedAt,
            UpdatedAt: chirp.UpdatedAt,
            UserID: chirp.ID, 
            Body: chirp.Body,
        })
    }

    jsonResponse, err := json.Marshal(chirpsFormatted)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(fmt.Appendf([]byte{}, `{"error": %s}`, err))
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}


func (cfg *apiConfig) handleGetSingleChirp (w http.ResponseWriter, r *http.Request){
    chirpID, err := uuid.Parse(r.PathValue("chirpID"))
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write(fmt.Appendf([]byte{}, `{"error":"%s"}`, err))
    }

    chirp, err := cfg.DB.GetSingleChirp(r.Context(), chirpID)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write(fmt.Appendf([]byte{}, `{"error":"%s"}`, err))
    }

    chirpFormatted := Chirp{
        ID: chirp.ID, 
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        UserID: chirp.ID, 
        Body: chirp.Body,
    }
    fmt.Println(chirpFormatted)

    jsonResponse, err := json.Marshal(chirpFormatted)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write(fmt.Appendf([]byte{}, `{"error":"%s"}`, err))
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

