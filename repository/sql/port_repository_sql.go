package repository

import (
	"database/sql"
	"fmt"

	"github.com/dscabral/ports/domain"
	"github.com/lib/pq"
)

type PortRepository struct {
	DB *sql.DB
}

func NewPortRepository(dbFile string) *PortRepository {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	return &PortRepository{
		DB: db,
	}
}

func (r *PortRepository) Init() error {
	_, err := r.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ports (
			id TEXT PRIMARY KEY,
			name TEXT,
			code TEXT,
			city TEXT,
			country TEXT,
			alias TEXT[],
			regions TEXT[],
			coordinates TEXT[],
			province TEXT,
			timezone TEXT,
			unlocs TEXT[]
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create ports table %s", err)
	}
	return nil
}

func (r *PortRepository) Close() error {
	return r.DB.Close()
}

func (r *PortRepository) InsertOrUpdatePort(port domain.Port) error {
	stmt, err := r.DB.Prepare(`INSERT OR REPLACE INTO ports (id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %w", err)
	}
	defer stmt.Close()

	// aliasArray := pq.Array(port.Alias)
	// regionsArray := pq.Array(port.Regions)
	coordinates := make([]string, len(port.Coordinates))
	for i, coord := range port.Coordinates {
		coordinates[i] = fmt.Sprintf("%f", coord)
	}
	// coordinatesArray := pq.Array(coordinates)

	_, err = stmt.Exec(port.ID, port.Name, port.City, port.Country, pq.Array(port.Alias), pq.Array(port.Regions), pq.Array(coordinates), port.Province, port.Timezone, pq.Array(port.Unlocs), port.Code)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %w", err)
	}

	return nil
}
