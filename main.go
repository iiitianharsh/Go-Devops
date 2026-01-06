package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	AppName string `json:"app_name"`
	Port    string `json:"port"`
	Env     string `json:"env"`
}

func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	return Config{
		AppName: "go-devops-app",
		Port:    port,
		Env:     env,
	}
}

func writeFile() {
	f, err := os.Create("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "application started")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("worker %d running\n", id)
	time.Sleep(2 * time.Second)
	log.Printf("worker %d done\n", id)
}

func healthHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg)
	}
}

func main() {
	cfg := loadConfig()
	writeFile()

	log.Printf("starting %s in %s mode\n", cfg.AppName, cfg.Env)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler(cfg))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("http server on :%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	wg.Wait()
	log.Println("clean exit")
}
