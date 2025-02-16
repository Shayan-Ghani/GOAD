package tagsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tagrequest "github.com/Shayan-Ghani/GOAD/pkg/request/tag"
)

func Handle(mux *http.ServeMux, s Service) {
	mux.HandleFunc("POST /tags", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)

			return
		}
		defer r.Body.Close()

		var tags tagrequest.Tag
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&tags); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)

			return
		}

		if err := s.Add(tags); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tags)

	})

	mux.HandleFunc("DELETE /tags/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")


		if err := s.Delete(tagrequest.Delete{
			Name: name,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	})



	mux.HandleFunc("POST /tags/item", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)

			return
		}
		defer r.Body.Close()

		var it tagrequest.BasePayload
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&it); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)

			return
		}

		if err := s.AddToItem(it); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		json.NewEncoder(w).Encode(it)

	})


	mux.HandleFunc("DELETE /tags/item/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var it = &tagrequest.BasePayload{}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(it); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)

			return
		}
		it.ItemID = id

		if err:= s.DeleteFromItem(*it); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	})


	mux.HandleFunc("GET /tags/item/{itemID}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")

		t, err := s.GetFromItems(tagrequest.Base{
			ItemID: id,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		json.NewEncoder(w).Encode(t)
	})

	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	// TODO: change all to query parametr 
	mux.HandleFunc("DELETE /tags/item/{itemID}/all", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("itemID")
		if err := s.DeleteAllFromItem(tagrequest.Base{
			ItemID: id,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	})
}