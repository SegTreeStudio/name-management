package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/SegTreeStudio/name-management/backend/internal/config"
    "github.com/SegTreeStudio/name-management/backend/internal/database"
    "github.com/SegTreeStudio/name-management/backend/internal/handler"
    "github.com/SegTreeStudio/name-management/backend/internal/repository"
)

func main() {

    // 1. Load configuration.
    cfg := config.LoadConfig()

    // 2. Connect to PostgreSQL.
    db, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	log.Println("Connected to PostgreSQL successfully")

    // 3. Create repository.
    nameRepo := repository.NewNameRepository(db)

    // 4. Create handler.
    nameHandler := handler.NewNameHandler(nameRepo)

    // 5. Register API routes.
    mux := http.NewServeMux()

    mux.HandleFunc("POST /api/names", nameHandler.CreateName)
    mux.HandleFunc("GET /api/names/exists", nameHandler.SearchName)
    mux.HandleFunc("GET /api/names", nameHandler.GetAllNames)

    mux.HandleFunc("GET /health", func(
        w http.ResponseWriter,
        r *http.Request,
    ) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        if err := db.Ping(ctx); err != nil {
            http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Server is running"))
    })

    // 6. Enable CORS.
    server := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: corsMiddleware(cfg.FrontendURL, mux),
    }

    log.Printf("Server running on http://localhost:%s", cfg.Port)

    if err := server.ListenAndServe(); err != nil &&
        err != http.ErrServerClosed {
        log.Fatal("Server error:", err)
    }
}

func corsMiddleware(
    allowedOrigin string,
    next http.Handler,
) http.Handler {

    return http.HandlerFunc(func(
        w http.ResponseWriter,
        r *http.Request,
    ) {

        origin := r.Header.Get("Origin")

        if origin == allowedOrigin {
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
            w.Header().Set("Vary", "Origin")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        }

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}