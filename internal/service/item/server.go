package itemsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	itemrequest "github.com/Shayan-Ghani/GOAD/pkg/request/item"
)


func Handle(mux *http.ServeMux, s *Service) {

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()


		var item itemrequest.Add
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&item); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received item: %+v\n", item)

		if err = s.Add(item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(item)
	})

	mux.HandleFunc("DELETE /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := s.Delete(itemrequest.Delete{
			ID: id,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(http.StatusOK)

	})

	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		res, err := s.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(res.Items)
	})

	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		res, err := s.GetSingle(itemrequest.Get{
			ID: id,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if res != nil{
			json.NewEncoder(w).Encode(res.Items)
		}
	})
	mux.HandleFunc("PUT /items", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var item itemrequest.Update
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&item); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err := s.Update(item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(http.StatusNoContent)
	})

	mux.HandleFunc("POST /items/done", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var item itemrequest.UpdateStatus
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&item); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err:= s.UpdateStatus(itemrequest.UpdateStatus{
			ID: item.ID,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(http.StatusCreated)
	})

	mux.HandleFunc("GET /items/done", func(w http.ResponseWriter, r *http.Request) {
		res, err := s.GetDone()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(res.Items)
	})
}
