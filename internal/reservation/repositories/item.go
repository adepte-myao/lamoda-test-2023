package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type PostgresItemRepository struct {
	db *sql.DB
}

func NewPostgresItem(db *sql.DB) *PostgresItemRepository {
	return &PostgresItemRepository{db: db}
}

func (repo PostgresItemRepository) GetAllAsMap(ctx context.Context) (map[domain.ItemID]domain.Item, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, length_meters, width_meters, height_meters, weight_kg FROM items`)
	if err != nil {
		return nil, fmt.Errorf("looking up in items table: %w", err)
	}

	items := make(map[domain.ItemID]domain.Item)
	for rows.Next() {
		var item domain.Item
		err = rows.Scan(&item.ID, &item.Name,
			&item.Size.LengthMeters, &item.Size.WidthMeters, &item.Size.HeightMeters, &item.WeightKilograms)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		items[item.ID] = item
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("after iterating over items rows: %w", err)
	}

	return items, nil
}
