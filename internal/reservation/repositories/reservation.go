package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/adepte-myao/lamoda-test-2023/internal/reservation/core/domain"
)

type PostgresReservationRepository struct {
	db *sql.DB
}

func NewPostgresReservation(db *sql.DB) *PostgresReservationRepository {
	return &PostgresReservationRepository{db: db}
}

func (repo PostgresReservationRepository) GetByID(ctx context.Context, id string) (domain.Reservation, error) {
	reservation := domain.Reservation{ID: id}

	err := repo.db.QueryRowContext(ctx,
		`SELECT destination_latitude, destination_longitude FROM reservations WHERE id = $1`, id,
	).Scan(&reservation.DestinationLocation.Latitude, &reservation.DestinationLocation.Longitude)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("looking up in reservation table: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx,
		`SELECT item_id, storehouse_id, items_count FROM reservation_items WHERE reservation_id = $1`, id)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("looking up in reservation_items table: %w", err)
	}

	for rows.Next() {
		var entry domain.ReserveEntry
		err = rows.Scan(&entry.ItemID, &entry.SourceStorehouseID, &entry.Count)
		if err != nil {
			return domain.Reservation{}, fmt.Errorf("scanning row: %w", err)
		}

		reservation.Entries = append(reservation.Entries, entry)
	}

	if rows.Err() != nil {
		return domain.Reservation{}, fmt.Errorf("after iterating over reservation_items rows: %w", err)
	}

	return reservation, nil
}

func (repo PostgresReservationRepository) Save(ctx context.Context, reservation domain.Reservation) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		// TODO
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO reservations (id, destination_latitude, destination_longitude) VALUES ($1, $2, $3)`,
		reservation.ID, reservation.DestinationLocation.Latitude, reservation.DestinationLocation.Longitude)
	if err != nil {
		return fmt.Errorf("inserting into reservations table: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO reservation_items (reservation_id, item_id, storehouse_id, items_count) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("preparing statement for reservation_items: %w", err)
	}

	var resultErr error
	for _, entry := range reservation.Entries {
		_, err = stmt.ExecContext(ctx, reservation.ID, entry.ItemID, entry.SourceStorehouseID, entry.Count)
		if err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("executing statement for reservation_items: %w", err))
		}
	}

	resultErr = errors.Join(resultErr, tx.Commit())
	if resultErr != nil {
		return fmt.Errorf("statement work and transaction commitment: %w", err)
	}

	return nil
}

func (repo PostgresReservationRepository) Update(ctx context.Context, reservation domain.Reservation) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		// TODO
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx,
		`UPDATE reservations SET destination_latitude = $2, destination_longitude = $3 WHERE id = $1`,
		reservation.ID, reservation.DestinationLocation.Latitude, reservation.DestinationLocation.Longitude)
	if err != nil {
		return fmt.Errorf("updating reservations table: %w", err)
	}

	// TODO: calculate changes instead of deleting-inserting all content
	_, err = tx.ExecContext(ctx,
		`DELETE FROM reservation_items WHERE reservation_id = $1`, reservation.ID)
	if err != nil {
		return fmt.Errorf("deleting associated reservation_items: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO reservation_items (reservation_id, item_id, storehouse_id, items_count) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("preparing statement for reservation_items: %w", err)
	}

	var resultErr error
	for _, entry := range reservation.Entries {
		_, err = stmt.ExecContext(ctx, reservation.ID, entry.ItemID, entry.SourceStorehouseID, entry.Count)
		if err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("executing statement for reservation_items: %w", err))
		}
	}

	resultErr = errors.Join(resultErr, tx.Commit())
	if resultErr != nil {
		return fmt.Errorf("statement work and transaction commitment: %w", err)
	}

	return nil
}

func (repo PostgresReservationRepository) Delete(ctx context.Context, id string) error {
	_, err := repo.db.ExecContext(ctx,
		`DELETE FROM reservation_items WHERE reservation_id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting associated reservation_items: %w", err)
	}

	_, err = repo.db.ExecContext(ctx,
		`DELETE FROM reservations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting reservation: %w", err)
	}

	return nil
}
