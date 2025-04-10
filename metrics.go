package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	// "strconv"
)

func (cfg *apiConfig) handleMetrics (w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-Type", "text/html")
    templ, err := template.ParseFiles("./admin_metrics.html")

    if err != nil {
        fmt.Println(err)
        w.Write([]byte("This did not work"))
    }
    // data := map[string]string{"HITS": strconv.Itoa(int(cfg.fileserverHits.Load()))} 
    err = templ.Execute(w, map[string]string{"HITS": strconv.Itoa(int(cfg.fileserverHits.Load()))})
    if err != nil {
        fmt.Println(err)
        w.Write([]byte("This (second thing) did not work"))
    }
    // w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileserverHits.Load()))))
}

func (cfg *apiConfig) handleReset (w http.ResponseWriter, r *http.Request){
    cfg.fileserverHits.Store(0)
}


func (cfg *apiConfig) middlewareMetricsInc (next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        // fmt.Println("Incrementing cfg.fileserverHits")
        cfg.fileserverHits.Store(cfg.fileserverHits.Add(1))
        next.ServeHTTP(w, r)
        // fmt.Println("THis comes after the next handler")
    })
}


