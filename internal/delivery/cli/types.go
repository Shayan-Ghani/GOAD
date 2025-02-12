package cli

import (
    "github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
    "github.com/Shayan-Ghani/GOAD/internal/model"
)

type CliRequest struct {
    Resource string
    Command  string
    Flags    *cmdflag.Flags
}

type CliResponse struct {
    Err   string       `json:"error,omitempty"`
    Items []model.Item `json:"items"`
}