package handler

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5/pgconn"

    "github.com/SegTreeStudio/name-management/backend/internal/model"
    "github.com/SegTreeStudio/name-management/backend/internal/repository"
)

type NameHandler struct {
    Repo *repository.NameRepository
}

func NewNameHandler(repo *repository.NameRepository) *NameHandler {
    return &NameHandler{
        Repo: repo,
    }
}

func writeJSON(
    w http.ResponseWriter,
    status int,
    data any,
) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}

// POST /api/names
func (h *NameHandler) CreateName(
    w http.ResponseWriter,
    r *http.Request,
) {

    var req model.CreateNameRequest

    r.Body = http.MaxBytesReader(w, r.Body, 4096)

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
            Error: "Invalid JSON request",
        })
        return
    }

    name := repository.NormalizeName(req.Name)

    if name == "" {
        writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
            Error: "Name cannot be empty",
        })
        return
    }

    if len([]rune(name)) > 100 {
        writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
            Error: "Name cannot exceed 100 characters",
        })
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    err := h.Repo.Create(ctx, name)

    if err != nil {

        var pgErr *pgconn.PgError

        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            writeJSON(w, http.StatusConflict, model.ErrorResponse{
                Error: "This name already exists",
            })
            return
        }

        writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
            Error: "Unable to save name",
        })
        return
    }

    writeJSON(w, http.StatusCreated, model.NameResponse{
        Message: "Name added successfully",
        Name:    name,
    })
}

// GET /api/names/exists?name=Asif
func (h *NameHandler) SearchName(
    w http.ResponseWriter,
    r *http.Request,
) {

    name := repository.NormalizeName(r.URL.Query().Get("name"))

    if name == "" {
        writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
            Error: "Please provide a name",
        })
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    exists, err := h.Repo.Exists(ctx, name)

    if err != nil {
        writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
            Error: "Unable to search name",
        })
        return
    }

    writeJSON(w, http.StatusOK, model.SearchNameResponse{
        Name:   name,
        Exists: exists,
    })
}

// GET /api/names
func (h *NameHandler) GetAllNames(
    w http.ResponseWriter,
    r *http.Request,
) {

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    names, err := h.Repo.GetAll(ctx)

    if err != nil {
        writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{
            Error: "Unable to retrieve names",
        })
        return
    }

    writeJSON(w, http.StatusOK, names)
}