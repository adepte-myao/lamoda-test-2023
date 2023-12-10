package repositories

import (
	"database/sql"
)

type PostgresReservationRepository struct {
	db *sql.DB
}
