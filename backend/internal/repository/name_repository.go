package repository

import (
    "context"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/SegTreeStudio/name-management/backend/internal/model"
)

type NameRepository struct {
    DB *pgxpool.Pool
}

func NewNameRepository(db *pgxpool.Pool) *NameRepository {
    return &NameRepository{
        DB: db,
    }
}

// Insert a new name into PostgreSQL.
func (r *NameRepository) Create(
    ctx context.Context,
    name string,
) error {

    query := `
        INSERT INTO names (name)
        VALUES ($1)
    `

    _, err := r.DB.Exec(ctx, query, name)

    return err
}

// Check whether a name exists.
func (r *NameRepository) Exists(
    ctx context.Context,
    name string,
) (bool, error) {

    query := `
        SELECT EXISTS (
            SELECT 1
            FROM names
            WHERE LOWER(name) = LOWER($1)
        )
    `

    var exists bool

    err := r.DB.QueryRow(ctx, query, name).Scan(&exists)

    if err != nil {
        return false, err
    }

    return exists, nil
}

// Return all names.
func (r *NameRepository) GetAll(
    ctx context.Context,
) ([]model.Name, error) {

    query := `
        SELECT id, name, created_at
        FROM names
        ORDER BY id DESC
    `

    rows, err := r.DB.Query(ctx, query)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    names := make([]model.Name, 0)

    for rows.Next() {

        var item model.Name

        err := rows.Scan(
            &item.ID,
            &item.Name,
            &item.CreatedAt,
        )

        if err != nil {
            return nil, err
        }

        names = append(names, item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return names, nil
}

// Normalize user input.
func NormalizeName(name string) string {
    return strings.TrimSpace(name)
}