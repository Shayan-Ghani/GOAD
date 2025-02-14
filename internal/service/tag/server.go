package tagsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	itemtagrequest "github.com/Shayan-Ghani/GOAD/pkg/request/tag"
)

func Handle(mux *http.ServeMux, s Service) {
	mux.HandleFunc("POST /tags", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var tags itemtagrequest.Tag
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&tags); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.Add(tags); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tags)

	})
	mux.HandleFunc("POST /tags/item", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var it itemtagrequest.BasePayload
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&it); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.AddToItem(it); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(it)

	})
	mux.HandleFunc("DELETE /tags/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var it = &itemtagrequest.BasePayload{}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(it); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}
		it.ItemID = id

		s.DeleteFromItem(*it)

	})
	mux.HandleFunc("GET /tags/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")

		t, err := s.GetFromItems(itemtagrequest.Base{
			ItemID: id,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(t)
	})

	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	mux.HandleFunc("DELETE /tags/{itemID}/all", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")
		if err := s.DeleteAllFromItem(itemtagrequest.Base{
			ItemID: id,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		s.DeleteAllFromItem(itemtagrequest.Base{
			ItemID: id,
		})

		w.WriteHeader(http.StatusNoContent)
	})
}