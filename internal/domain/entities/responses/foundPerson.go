package responses

import "github.com/agl/fio/internal/domain/entities"

type FoundPerson struct {
	Message string `json:"message"`
	Data    entities.Person `json:"data"`
}