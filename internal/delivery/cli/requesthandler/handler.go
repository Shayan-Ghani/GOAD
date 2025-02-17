package requesthandler

import (
	"fmt"
	"net/http"
	"time"

)

type Handler struct {
	client     *http.Client
	ItemSvcUrl string
	TagSvcUrl  string
}


func (cr *CliResponse) Error() string {
	return cr.Err
}

func NewHandler(ItemSvcUrl string, TagSvcUrl string) Handler {
	return Handler{
		ItemSvcUrl: ItemSvcUrl,
		TagSvcUrl:  TagSvcUrl,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h Handler) HandleRequest(request CliRequest) (*CliResponse, error) {

	switch request.Resource {
	case "item":
		return h.handleItem(request)
	case "tag":
		return nil, h.handleTag(request)
	default:
		return nil, fmt.Errorf("unkown resource %v", request.Resource)
	}

}