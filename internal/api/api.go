package api

import (
	"encoding/json"
	"net/http"

	"github.com/SyedAsadK/govecDB/internal/db"
)

type InsertPayload struct {
	ID     string    `json:"id"`
	Vector []float64 `json:"vector"`
}
type SearchPayload struct {
	Vector []float64 `json:"vector"`
	K      int       `json:"k"`
}
type SearchResponse struct {
	ID    string    `json:"id"`
	Score float64   `json:"score"`
	Data  []float64 `json:"data"`
}

type Controller struct {
	Store *db.VectorStore
}

func (c *Controller) HandleInsert(w http.ResponseWriter, r *http.Request) {
	var payload InsertPayload
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&payload)
	if err != nil {
		http.Error(w, "error in handleInsert : bad json", http.StatusInternalServerError)
		return
	}
	dp := db.NewVectorPoint(payload.ID, payload.Vector)
	err = c.Store.Insert(dp)
	if err != nil {
		http.Error(w, "error in handleInsert : storing data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var payload SearchPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "error in handleSearch: bad json", http.StatusInternalServerError)
		return
	}
	if payload.K <= 0 {
		payload.K = 1
	}
	res, err := c.Store.Search(payload.Vector, payload.K)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	apiResults := make([]SearchResponse, 0)
	for _, v := range res {
		apiResults = append(apiResults, SearchResponse{
			ID:    v.Vector.Id,
			Score: v.Score,
			Data:  v.Vector.Vector,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apiResults)
	if err != nil {
		http.Error(w, "error in handleSearch: encoding error", http.StatusInternalServerError)
		return
	}
}
