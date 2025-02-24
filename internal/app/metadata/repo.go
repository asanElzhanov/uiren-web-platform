package metadata

import (
	"context"
	"encoding/json"
	"uiren/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type metadataRepository struct {
	db *pgxpool.Pool
}

func NewMetadataRepository(db *pgxpool.Pool) *metadataRepository {
	return &metadataRepository{
		db: db,
	}
}

func (r *metadataRepository) getMetadata(ctx context.Context, dictID int) ([]Metadata, error) {
	var (
		query = `
		SELECT keys FROM dictionaries where id = $1;
		`
	)
	row := r.db.QueryRow(ctx, query, dictID)

	var jsonBody json.RawMessage
	var metadataList []Metadata
	if err := row.Scan(&jsonBody); err != nil {
		logger.Error("Ошибка Scan: ", err)
		return nil, err
	}
	if err := json.Unmarshal(jsonBody, &metadataList); err != nil {
		logger.Error("Ошибка Unmarshal: ", err)
		return nil, err
	}
	return metadataList, nil
}
