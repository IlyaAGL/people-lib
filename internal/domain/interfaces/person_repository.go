package interfaces

import (
	"github.com/agl/fio/internal/domain/entities"
)

type PersonRepository interface {
	DeletePersonByID(id int) error
	GetPersonByID(id int) (entities.Person, error)
	GetPeopleByFilter(filter entities.Person, page, limit string) ([]entities.Person, error)
	CreatePerson(person entities.Person) (int, error)
	UpdatePersonByID(id int, p entities.Person) error
}
