package model

import "time"

type Name struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateNameRequest struct {
    Name string `json:"name"`
}

type NameResponse struct {
    Message string `json:"message"`
    Name    string `json:"name"`
}

type SearchNameResponse struct {
    Name   string `json:"name"`
    Exists bool   `json:"exists"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}