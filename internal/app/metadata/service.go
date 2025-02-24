package metadata

import (
	"context"
	"errors"
	"fmt"
	"uiren/internal/app/data"
)

type repository interface {
	getMetadata(ctx context.Context, dictId int) ([]Metadata, error)
}

type MetadataService struct {
	repo repository
}

func NewMetadataService(repo repository) *MetadataService {
	return &MetadataService{
		repo: repo,
	}
}

func (s MetadataService) GetMetadataTypes(ctx context.Context, dictId int) (map[string]string, error) {
	requisites, err := s.repo.getMetadata(ctx, dictId)
	if err != nil {
		return nil, err
	}
	metadataTypes := make(map[string]string)
	metadataTypes["Id"] = "number"
	metadataTypes["Parent"] = "number"
	metadataTypes["IsMarkedToDelete"] = "boolean"
	metadataTypes["created_at"] = "number"
	metadataTypes["deleted_at"] = "number"
	for _, metadataEntity := range requisites {
		dataType, ok := metadataEntity.Data.DataType["type"].(string)
		if !ok {
			return nil, errors.New("Ошибка, type в метаданных не строка")
		}
		metadataTypes[metadataEntity.Data.RequisiteName] = dataType
	}
	return metadataTypes, nil
}

func InsertData(data data.DictionaryData) {
	fmt.Println(data)
	fmt.Println("Загружается Бекарысом в бд")
}

func (s MetadataService) GetMetadataSortOrder(ctx context.Context, dictId int) (map[string]int, error) {
	requisites, err := s.repo.getMetadata(ctx, dictId)
	if err != nil {
		return nil, err
	}

	metadataSortOrder := make(map[string]int)
	metadataSortOrder["Id"] = 0
	metadataSortOrder["Parent"] = 1
	metadataSortOrder["IsMarkedToDelete"] = 2
	for _, requisite := range requisites {
		metadataSortOrder[requisite.Data.RequisiteName] = int(requisite.Data.SortOrder) + 2
	}

	return metadataSortOrder, nil
}
