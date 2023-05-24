package mocks

import (
	"github.com/dscabral/ports/src/domain"
)

type MockPortRepository struct {
	PortData map[string]domain.Port
}

func (m *MockPortRepository) InsertOrUpdatePort(port domain.Port) error {
	m.PortData[port.ID] = port
	return nil
}

func NewMockPortRepository() *MockPortRepository {
	return &MockPortRepository{
		PortData: make(map[string]domain.Port),
	}
}
