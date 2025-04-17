package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
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
    mux.HandleFunc("POST /api/users", apiCnfg.handleCreateUser)
    mux.HandleFunc("POST /api/chirps", apiCnfg.handleCreateChirp)

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

