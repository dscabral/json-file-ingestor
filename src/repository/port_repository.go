package repository

import "github.com/dscabral/ports/src/domain"

type PortRepository interface {
	InsertOrUpdatePort(port domain.Port) error
}
