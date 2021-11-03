package db

/*
type RDBConnector interface {
	Add() int
	AddBatch() int
	Update() int
	Delete() int
	Search() int
}
*/

type localMemoryConnector struct {
	// key: table name; id: uuid
	storage map[string]map[string]interface{}
}

func (connector *localMemoryConnector) init() {
	connector.storage = make(map[string]map[string]interface{})
}

func (connector *localMemoryConnector) Add(config RDBAddConfig) int {
	// todo: add into storage
}
