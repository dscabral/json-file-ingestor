package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

type Port struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
}

const (
	DBFile = "ports.db"
)

func main() {
	db, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ports (
			id TEXT PRIMARY KEY,
			name TEXT,
			code TEXT,
			city TEXT,
			country TEXT,
			alias TEXT,
			regions TEXT,
			coordinates TEXT,
			province TEXT,
			timezone TEXT,
			unlocs TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("example.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var portMap map[string]Port
		if err := decoder.Decode(&portMap); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		for key, port := range portMap {
			port.ID = key
			insertOrUpdatePort(db, port)
		}
	}

	http.HandleFunc("/ports", createOrUpdatePortHandler)
	http.HandleFunc("/ports/all", getAllPortsHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	// Clean up the database file after shutting down
	if err := os.Remove(DBFile); err != nil {
		log.Println(err)
	}
}

func insertOrUpdatePort(db *sql.DB, port Port) {
	stmt, err := db.Prepare(`
		INSERT INTO ports (id, name, code, city, country, alias, regions, coordinates, province, timezone, unlocs)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			code = excluded.code,
			city = excluded.city,
			country = excluded.country,
			alias = excluded.alias,
			regions = excluded.regions,
			coordinates = excluded.coordinates,
			province = excluded.province,
			timezone = excluded.timezone,
			unlocs = excluded.unlocs
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	aliasJSON, _ := json.Marshal(port.Alias)
	regionsJSON, _ := json.Marshal(port.Regions)
	coordinatesJSON, _ := json.Marshal(port.Coordinates)
	unlocsJSON, _ := json.Marshal(port.Unlocs)

	_, err = stmt.Exec(
		port.ID,
		port.Name,
		port.Code,
		port.City,
		port.Country,
		string(aliasJSON),
		string(regionsJSON),
		string(coordinatesJSON),
		port.Province,
		port.Timezone,
		string(unlocsJSON),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func createOrUpdatePortHandler(w http.ResponseWriter, r *http.Request) {
	var port Port
	err := json.NewDecoder(r.Body).Decode(&port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insertOrUpdatePort(db, port)
	fmt.Fprintln(w, "Port created or updated successfully")
}

func getAllPortsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ports")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ports []Port

	for rows.Next() {
		var port Port
		var aliasJSON, regionsJSON, coordinatesJSON, unlocsJSON string
		err := rows.Scan(
			&port.ID,
			&port.Name,
			&port.Code,
			&port.City,
			&port.Country,
			&aliasJSON,
			&regionsJSON,
			&coordinatesJSON,
			&port.Province,
			&port.Timezone,
			&unlocsJSON,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal([]byte(aliasJSON), &port.Alias)
		if err != nil {
			http.Error(w, "Failed to parse port alias", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal([]byte(regionsJSON), &port.Regions)
		if err != nil {
			http.Error(w, "Failed to parse port regions", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal([]byte(coordinatesJSON), &port.Coordinates)
		if err != nil {
			http.Error(w, "Failed to parse port coordinates", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal([]byte(unlocsJSON), &port.Unlocs)
		if err != nil {
			http.Error(w, "Failed to parse port UNLOCs", http.StatusInternalServerError)
			return
		}

		ports = append(ports, port)
	}

	portsJSON, err := json.Marshal(ports)
	if err != nil {
		http.Error(w, "Failed to serialize port data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(portsJSON)
}
