package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
    "log"
	"github.com/MalikL2005/http_server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
    DB *database.Queries
    fileserverHits atomic.Int32
}


func main(){
    godotenv.Load()

    db_url := os.Getenv("DB_URL")
    if db_url == "" {
        log.Fatal("Could not get DB_URL")
    }
    db, err := sql.Open("postgres", db_url)
    if err != nil {
        log.Fatal(fmt.Sprintf("Database error: %s", err))
    }
    db_queries := database.New(db)

    const filePathRoot = "."
    const port = "8080"

    mux := http.NewServeMux()
    apiCnfg := apiConfig{DB: db_queries}
    apiCnfg.fileserverHits.Store(0)
    mux.Handle("/app/", http.StripPrefix("/app", apiCnfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
    mux.HandleFunc("GET /api/healthz", handleHealthz)
    mux.HandleFunc("GET /admin/metrics", apiCnfg.handleMetrics)
    mux.HandleFunc("POST /admin/reset", apiCnfg.handleReset)
    mux.HandleFunc("POST /api/validate_chirp", handleValidation)
    mux.HandleFunc("POST /api/users", apiCnfg.handleCreateUser)

    server := http.Server{Handler: mux, Addr: ":" + port};
    fmt.Println("Serving on port ", port)
    err = server.ListenAndServe()
    if err != nil {
        log.Fatal("Errrorrororrorr")
    }

}


func handleHealthz (writer http.ResponseWriter, res *http.Request){
    writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
    writer.WriteHeader(http.StatusOK)
    writer.Write([]byte(http.StatusText(http.StatusOK)))
}

var unallowed_words []string = []string{
    "kerfuffle",
    "sharbert",
    "fornax",
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
    writer.Write([]byte(fmt.Sprintf(`{"cleaned_body": "%s"}`, val.Body)))
}

