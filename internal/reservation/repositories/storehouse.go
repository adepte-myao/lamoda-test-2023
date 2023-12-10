package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

// TODO: transactions

type PostgresStorehouseRepository struct {
	db *sql.DB
}

func NewPostgresStorehouse(db *sql.DB) *PostgresStorehouseRepository {
	return &PostgresStorehouseRepository{db: db}
}

func (repo PostgresStorehouseRepository) GetItemsByID(ctx context.Context, id domain.StoreHouseID) (map[domain.ItemID]domain.ItemData, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT item_id, items_count FROM storehouses_items WHERE storehouse_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("looking up in storehouses_items table: %w", err)
	}

	// TODO: should add full info about items, not only their IDs?
	unreserved := make(map[domain.ItemID]domain.ItemData)
	for rows.Next() {
		var itemData domain.ItemData
		err = rows.Scan(&itemData.Item.ID, &itemData.Count)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		unreserved[itemData.Item.ID] = itemData
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("after iterating over storehouses_items rows: %w", err)
	}

	return unreserved, nil
}

func (repo PostgresStorehouseRepository) GetAllAsMap(ctx context.Context) (map[domain.StoreHouseID]domain.StoreHouse, error) {
	rows, err := repo.db.QueryContext(ctx, `SELECT id, name, latitude, longitude FROM storehouses`)
	if err != nil {
		return nil, fmt.Errorf("looking up in storehouses table: %w", err)
	}

	storehouses := make(map[domain.StoreHouseID]domain.StoreHouse)
	for rows.Next() {
		var storehouse domain.StoreHouse
		err = rows.Scan(&storehouse.ID, &storehouse.Name, &storehouse.Location.Latitude, &storehouse.Location.Longitude)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		storehouses[storehouse.ID] = storehouse
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("after iterating over storehouses rows: %w", err)
	}

	for id, storehouse := range storehouses {
		itemsData, err := repo.GetItemsByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("subquery for storehouses_items: %w", err)
		}

		storehouse.ItemsData = itemsData
		storehouses[id] = storehouse
	}

	return storehouses, nil
}

func (repo PostgresStorehouseRepository) UpdateAll(ctx context.Context, storehouses map[domain.StoreHouseID]domain.StoreHouse) error {
	//TODO: make better query instead of multiple delete-insert. Maybe domain logic is the one that should be changed
	deleteStmt, err := repo.db.PrepareContext(ctx, `DELETE FROM storehouses_items WHERE storehouse_id = $1`)
	if err != nil {
		return fmt.Errorf("preparing deletion query: %w", err)
	}

	insertStmt, err := repo.db.PrepareContext(ctx,
		`INSERT INTO storehouses_items (storehouse_id, item_id, items_count) VALUES ($1, $2, $3)`)
	if err != nil {
		return fmt.Errorf("preparing insertion query: %w", err)
	}

	for _, storehouse := range storehouses {
		_, err = deleteStmt.ExecContext(ctx, storehouse.ID)
		if err != nil {
			return fmt.Errorf("deleting: %w", err)
		}

		for itemID, itemData := range storehouse.ItemsData {
			_, err = insertStmt.ExecContext(ctx, storehouse.ID, itemID, itemData.Count)
			if err != nil {
				return fmt.Errorf("inserting: %w", err)
			}
		}
	}

	return nil
}
