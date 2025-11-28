package db

import (
	"encoding/gob"
	"errors"
	"os"
)

func (vs *VectorStore) Save(filename string) error {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("os: error creating file")
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(vs.Data)
	if err != nil {
		return errors.New("encode: error in encoding db")
	}
	return nil
}

func Load(filename string) (*VectorStore, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("os: db file could not be opened")
	}
	var data map[string]*VectorPoint
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, errors.New("decode: error in decoding db")
	}
	vs := NewVectorStore()
	vs.Data = data
	for _, v := range data {
		vs.Dimension = len(v.Vector)
		break
	}
	return vs, nil
}
