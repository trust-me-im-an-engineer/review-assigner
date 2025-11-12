package service

import "review-assigner/internal/storage"

// Service is considered to be a core layer of witch only one could exist,
// therefore it doesn't use interface
type Service struct {
	storage storage.Storage
}

func NewService(storage *storage.Storage) *Service {
	return &Service{storage: storage}
}
