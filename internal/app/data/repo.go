package data

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dataRepository struct {
	db *pgxpool.Pool
}

func NewDataRepository(db *pgxpool.Pool) *dataRepository {
	return &dataRepository{
		db: db,
	}
}

func (r *dataRepository) getDictionaryData(ctx context.Context, dictID int) ([]DictionaryData, error) {
	var (
		query = `
		SELECT data from dictionaries_data
		where dictionary_id = $1 AND deleted_at IS NULL;
		`
		response []DictionaryData
	)

	rows, err := r.db.Query(ctx, query, dictID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var jsonBody json.RawMessage
		if err := rows.Scan(&jsonBody); err != nil {
			return nil, err
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal(jsonBody, &jsonData); err != nil {
			return nil, err
		}

		dictData, ok := jsonData["data"].(map[string]interface{})
		if !ok {
			return nil, errors.New("corrupted data: data")
		}
		dictSystemData, ok := jsonData["system_data"].(map[string]interface{})
		if !ok {
			return nil, errors.New("corrupted data: systemData")
		}

		var DictionaryData DictionaryData
		DictionaryData.Data = dictData
		if createdAtFloat, ok := dictSystemData["created_at"].(float64); ok {
			DictionaryData.SystemData.Created_at = createdAtFloat
		} else {
			return nil, errors.New("bad request: invalid 'created_at'")
		}

		if deletedAtFloat, ok := dictSystemData["deleted_at"].(float64); ok {
			DictionaryData.SystemData.Deleted_at = deletedAtFloat
		}

		response = append(response, DictionaryData)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return response, nil
}

func (r *dataRepository) insertData(ctx context.Context, entityID int, dictData DictionaryData) (float64, error) {
	var (
		query = `
		INSERT INTO dictionaries_data(item_id, dictionary_id, parent_id, child_count, created_by, data) VALUES ($1, $2, $3, $4, $5, $6);
		`
	)

	jsonBody, err := json.Marshal(dictData)
	if err != nil {
		return -1, err
	}
	_, err = r.db.Exec(ctx, query, dictData.Data["Id"], entityID, 0, 0, 0, jsonBody)
	if err != nil {
		return -1, err
	}

	return dictData.Data["Id"].(float64), nil
}
