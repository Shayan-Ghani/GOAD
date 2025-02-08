package tcphandler

import (
	"github.com/Shayan-Ghani/GOAD/internal/delivery/command/cmdflag"
	"github.com/Shayan-Ghani/GOAD/internal/model"
	"github.com/Shayan-Ghani/GOAD/internal/repository"
)

type RequestHandler struct {
    repo repository.Repository
}


type CliRequest struct {
	Resource string
	Command  string
	Flags    *cmdflag.Flags
}

type CliResponse struct {
	Err   string       `json:"error,omitempty"`
	Items []model.Item `json:"items"`
}

func (cr *CliResponse) Error() string {
	return cr.Err
}

func NewRequestHandler(repo repository.Repository) *RequestHandler {
	return &RequestHandler{
		repo: repo,
	}
}

func (h *RequestHandler) HandleRequest(request CliRequest) CliResponse {
    if err := request.Flags.CheckCommandFlags(); err != nil {
        return CliResponse{
            Err: err.Error(),
        }
    }

    switch request.Resource {
    case "item":
        return h.handleItem(request)
    case "tag":
        return CliResponse{
            Err: h.handleTag(request).Error(),
        }
    default:
        return CliResponse{
            Err: "unknown resource",
        }
    }
}