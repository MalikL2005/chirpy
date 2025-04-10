package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func (cfg *apiConfig) handleMetrics (w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-Type", "text/html")
    templ, err := template.ParseFiles("./admin_metrics.html")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Parsing of template failed"))
    }
    err = templ.Execute(w, map[string]string{"HITS": strconv.Itoa(int(cfg.fileserverHits.Load()))})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Template execution failed"))
    }
    // w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileserverHits.Load()))))
}

func (cfg *apiConfig) handleReset (w http.ResponseWriter, r *http.Request){
    platform := os.Getenv("PLATFORM")
    if platform != "dev"{
        w.WriteHeader(http.StatusForbidden)
        return
    }
    cfg.fileserverHits.Store(0)
    err := cfg.DB.DeleteAllUsers(context.Background())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}


func (cfg *apiConfig) middlewareMetricsInc (next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        cfg.fileserverHits.Store(cfg.fileserverHits.Add(1))
        next.ServeHTTP(w, r)
    })
}


