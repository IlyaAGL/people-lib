package interfaces

import "github.com/agl/fio/internal/domain/entities"

type PersonService interface {
	DeletePersonByID(idStr string) error
	GetPersonByID(idStr string) (entities.Person, error)
	GetPeopleByFilter(name, surname, ageStr, gender, nationality, page, limit string, patronymic *string) ([]entities.Person, error)
	CreatePerson(person entities.Person) (int, error)
	UpdatePersonByID(person entities.Person, idStr string) error
}
