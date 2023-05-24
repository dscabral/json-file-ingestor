// service/port_service.go
package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dscabral/ports/src/domain"
	"github.com/dscabral/ports/src/repository"
)

type PortService struct {
	PortRepository repository.PortRepository
}

func NewPortService(portRepository repository.PortRepository) *PortService {
	return &PortService{
		PortRepository: portRepository,
	}
}

func (s *PortService) SaveOrUpdatePortFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var portMap map[string]domain.Port
		if err := decoder.Decode(&portMap); err != nil {
			if err.Error() == "EOF" {
				fmt.Println("File reading completed.")
				break
			}
			log.Fatal(err)
		}

		for key, port := range portMap {
			port.ID = key
			err := s.PortRepository.InsertOrUpdatePort(port)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
