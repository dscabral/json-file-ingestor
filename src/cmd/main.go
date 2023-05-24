package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	repository "github.com/dscabral/ports/src/repository/sql"
	"github.com/dscabral/ports/src/service"
	_ "modernc.org/sqlite"
)

// type Port struct {
// 	ID          string    `json:"id"`
// 	Name        string    `json:"name"`
// 	Code        string    `json:"code"`
// 	City        string    `json:"city"`
// 	Country     string    `json:"country"`
// 	Alias       []string  `json:"alias"`
// 	Regions     []string  `json:"regions"`
// 	Coordinates []float64 `json:"coordinates"`
// 	Province    string    `json:"province"`
// 	Timezone    string    `json:"timezone"`
// 	Unlocs      []string  `json:"unlocs"`
// }

const (
	DBFile = "ports.db"
)

func main() {
	portsRepository := repository.NewPortRepository(DBFile)
	defer func() {
		fmt.Println("shutting down database connection")
		if err := portsRepository.Close(); err != nil {
			log.Fatalf("failed to shutdown database connection: %v", err)
		}
	}()

	portsRepository.Init()

	portService := service.NewPortService(portsRepository)
	path := "ports.json"
	err := portService.SaveOrUpdatePortFromFile(path)
	if err != nil {
		log.Fatalf("failed to import and save the ports: %v", err)
	}

	// db, err := sql.Open("sqlite3", DBFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS ports (
	// 		id TEXT PRIMARY KEY,
	// 		name TEXT,
	// 		code TEXT,
	// 		city TEXT,
	// 		country TEXT,
	// 		alias TEXT,
	// 		regions TEXT,
	// 		coordinates TEXT,
	// 		province TEXT,
	// 		timezone TEXT,
	// 		unlocs TEXT
	// 	)
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// file, err := os.Open("ports.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// decoder := json.NewDecoder(file)
	// for {
	// 	var portMap map[string]Port
	// 	if err := decoder.Decode(&portMap); err != nil {
	// 		if err.Error() == "EOF" {
	// 			fmt.Println("File reading completed.")
	// 			break
	// 		}
	// 		log.Fatal(err)
	// 	}

	// 	for key, port := range portMap {
	// 		port.ID = key
	// 		insertOrUpdatePort(db, port)
	// 	}
	// }

	// http.HandleFunc("/ports", createOrUpdatePortHandler)
	// http.HandleFunc("/ports/all", getAllPortsHandler)
	// go func() {
	// 	log.Fatal(http.ListenAndServe(":8081", nil))
	// }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down...", sig)

	// Clean up the database file after shutting down
	if err := os.Remove(DBFile); err != nil {
		log.Println(err)
	}
}

// func insertOrUpdatePort(db *sql.DB, port Port) {
// 	stmt, err := db.Prepare(`
// 		INSERT INTO ports (id, name, code, city, country, alias, regions, coordinates, province, timezone, unlocs)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 		ON CONFLICT (id) DO UPDATE SET
// 			name = excluded.name,
// 			code = excluded.code,
// 			city = excluded.city,
// 			country = excluded.country,
// 			alias = excluded.alias,
// 			regions = excluded.regions,
// 			coordinates = excluded.coordinates,
// 			province = excluded.province,
// 			timezone = excluded.timezone,
// 			unlocs = excluded.unlocs
// 	`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()

// 	aliasJSON, _ := json.Marshal(port.Alias)
// 	regionsJSON, _ := json.Marshal(port.Regions)
// 	coordinatesJSON, _ := json.Marshal(port.Coordinates)
// 	unlocsJSON, _ := json.Marshal(port.Unlocs)

// 	_, err = stmt.Exec(
// 		port.ID,
// 		port.Name,
// 		port.Code,
// 		port.City,
// 		port.Country,
// 		string(aliasJSON),
// 		string(regionsJSON),
// 		string(coordinatesJSON),
// 		port.Province,
// 		port.Timezone,
// 		string(unlocsJSON),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func createOrUpdatePortHandler(w http.ResponseWriter, r *http.Request) {
// 	var port Port
// 	err := json.NewDecoder(r.Body).Decode(&port)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	db, err := sql.Open("sqlite3", DBFile)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer db.Close()

// 	insertOrUpdatePort(db, port)
// 	fmt.Fprintln(w, "Port created or updated successfully")
// }

// func getAllPortsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Retrieve ports from the repository
// 	ports, err := portRepository.GetAllPorts()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Serialize ports to JSON
// 	portsJSON, err := serializePorts(ports)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set response headers and write JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(portsJSON)
// }

// // serializePorts serializes the given ports to JSON
// func serializePorts(ports []domain.Port) ([]byte, error) {
// 	// Create a slice to hold the serialized port data
// 	serializedPorts := make([]map[string]interface{}, len(ports))

// 	// Serialize each port individually
// 	for i, port := range ports {
// 		serializedPort := map[string]interface{}{
// 			"id":          port.ID,
// 			"name":        port.Name,
// 			"code":        port.Code,
// 			"city":        port.City,
// 			"country":     port.Country,
// 			"alias":       port.Alias,
// 			"regions":     port.Regions,
// 			"coordinates": port.Coordinates,
// 			"province":    port.Province,
// 			"timezone":    port.Timezone,
// 			"unlocs":      port.Unlocs,
// 		}

// 		serializedPorts[i] = serializedPort
// 	}

// 	// Marshal the serialized ports to JSON
// 	return json.Marshal(serializedPorts)
// }
