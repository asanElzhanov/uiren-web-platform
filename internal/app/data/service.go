package data

import (
	"context"
)

type repository interface {
	getDictionaryData(ctx context.Context, dictID int) ([]DictionaryData, error)
	insertData(ctx context.Context, entityID int, dictData DictionaryData) (float64, error)
}

type DataService struct {
	repo repository
}

func NewDataService(repo repository) *DataService {
	return &DataService{
		repo: repo,
	}
}

func (s *DataService) GetDictionaryData(ctx context.Context, dictID int) ([]DictionaryData, error) {
	dataList, err := s.repo.getDictionaryData(ctx, dictID)
	if err != nil {
		return nil, err
	}
	return dataList, nil
}

func (s *DataService) InsertData(ctx context.Context, entityID int, dictData DictionaryData) (float64, error) {
	id, err := s.repo.insertData(ctx, entityID, dictData)
	if err != nil {
		return -1, err
	}
	return id, nil
}
