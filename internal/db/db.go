package db

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
)

type VectorPoint struct {
	Vector    []float64
	Id        string
	Magnitude float64
}
type VectorStore struct {
	Data      map[string]*VectorPoint
	Dimension int
	mu        sync.RWMutex
}

type SearchResult struct {
	Vector *VectorPoint
	Score  float64
}

func NewVectorStore() *VectorStore {
	return &VectorStore{
		Data: make(map[string]*VectorPoint),
	}
}

func (vs *VectorStore) Insert(v *VectorPoint) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	if _, ok := vs.Data[v.Id]; ok {
		return errors.New("insert: vector already in the vector store")
	}
	if len(v.Vector) == 0 {
		return errors.New("insert: Vector length should be greater than 0")
	}
	if vs.Dimension == 0 {
		vs.Dimension = len(v.Vector)
	} else if vs.Dimension != len(v.Vector) {
		return errors.New("insert: Vector should be of same length")
	}

	vs.Data[v.Id] = v
	return nil
}

func (vs *VectorStore) Delete(id string) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	if _, ok := vs.Data[id]; !ok {
		return errors.New("delete: id not found")
	}
	delete(vs.Data, id)
	return nil
}

func (vs *VectorStore) Search(qvec []float64, k int) ([]SearchResult, error) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	if vs.Dimension != len(qvec) {
		return nil, errors.New("search: dimension of input mismatched")
	}
	qVecPoint := NewVectorPoint("query", qvec)
	ranks := make([]SearchResult, 0, len(vs.Data))
	for _, i := range vs.Data {
		sim, _ := qVecPoint.CosineSimilarity(i)
		ranks = append(ranks, SearchResult{Vector: i, Score: sim})
	}
	sort.Slice(ranks, func(i, j int) bool { return ranks[i].Score > ranks[j].Score })
	if len(ranks) < k {
		k = len(ranks)
	}

	return ranks[:k], nil
}

func NewVectorPoint(id string, data []float64) *VectorPoint {
	return &VectorPoint{
		Id:        id,
		Vector:    data,
		Magnitude: vecMag(data),
	}
}

func (a *VectorPoint) CosineSimilarity(b *VectorPoint) (float64, error) {
	if a.Magnitude == 0 || b.Magnitude == 0 {
		return 0, errors.New("cannot calculate similarity with zero-Magnitude vector")
	}
	dotprod, err := dotProduct(a.Vector, b.Vector)
	if err != nil {
		return 0, errors.New("error in Dotproduct")
	}
	cosineSim := dotprod / (a.Magnitude * b.Magnitude)
	return cosineSim, nil
}

func dotProduct(v1, v2 []float64) (float64, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("DotProduct Error :Vectors were not of same size")
	}
	var sum float64
	for i := range v1 {
		sum += v1[i] * v2[i]
	}
	return sum, nil
}

func vecMag(v []float64) float64 {
	var sum float64
	for i := range v {
		sum += (v[i] * v[i])
	}
	return math.Sqrt(sum)
}
