package service_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dscabral/ports/domain"
	"github.com/dscabral/ports/repository/mocks"
	"github.com/dscabral/ports/service"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFile = "ports_test.db"
)

func TestSaveOrUpdatePortFromFile(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]domain.Port
		err  string
	}{
		{
			name: "SinglePort",
			data: map[string]domain.Port{
				"AEAJM": {
					ID:          "AEAJM",
					Name:        "Ajman",
					City:        "Ajman",
					Country:     "United Arab Emirates",
					Alias:       []string{},
					Regions:     []string{},
					Coordinates: []float64{55.5136433, 25.4052165},
					Province:    "Ajman",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{"AEAJM"},
					Code:        "52000",
				},
			},
			err: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary JSON file for testing
			tempFile, err := os.CreateTemp("", "ports_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFile.Name())

			// Encode test case port data to JSON and write to the temporary file
			jsonData, err := json.Marshal(tc.data)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := tempFile.Write(jsonData); err != nil {
				t.Fatal(err)
			}

			// Create an instance of PortRepository
			portRepo := mocks.NewMockPortRepository()

			// Create an instance of PortService with the PortRepository
			portService := service.NewPortService(portRepo)

			// Call the SaveOrUpdatePortFromFile method
			err = portService.SaveOrUpdatePortFromFile(tempFile.Name())

			if err != nil && tc.err == "" {
				t.Errorf("SaveOrUpdatePortFromFile returned an unexpected error: %v", err)
			} else if err == nil && tc.err != "" {
				t.Errorf("SaveOrUpdatePortFromFile did not return an expected error: %s", tc.err)
			} else if err != nil && err.Error() != tc.err {
				t.Errorf("SaveOrUpdatePortFromFile returned an unexpected error message: got '%v', want '%v'", err.Error(), tc.err)
			}
		})
	}
}
