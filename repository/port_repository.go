package repository

import "github.com/dscabral/ports/domain"

type PortRepository interface {
	InsertOrUpdatePort(port domain.Port) error
}
