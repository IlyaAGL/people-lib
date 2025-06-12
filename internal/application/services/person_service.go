package services

import (
	"errors"
	"strconv"

	"github.com/agl/fio/internal/domain/entities"
	"github.com/agl/fio/internal/domain/interfaces"
	. "github.com/agl/fio/pkg/logger"
)

type PersonService struct {
	repo interfaces.PersonRepository
}

func NewPersonService(repo interfaces.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (p *PersonService) UpdatePersonByID(person entities.Person, idStr string) error {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		Log.Info("Invalid ID format", "idStr", idStr, "error", err)

		return errors.New("invalid id format")
	}

	err = p.repo.UpdatePersonByID(id, person)
	if err != nil {
		Log.Info("Failed to update person", "id", id, "error", err, "person", person)

		return err
	}

	Log.Info("Person updated successfully", "id", id, "person", person)

	return nil
}

func (p *PersonService) CreatePerson(person entities.Person) (int, error) {
	id, err := p.repo.CreatePerson(person)

	if err != nil {
		Log.Info("Failed to create person", "error", err, "person", person)

		return 0, err
	}

	Log.Info("Person created successfully", "id", id, "person", person)

	return id, nil
}

func (p *PersonService) GetPersonByID(idStr string) (entities.Person, error) {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		Log.Info("Invalid ID format", "id", idStr, "error", err)

		return entities.Person{}, errors.New("invalid id format")
	}

	person, err := p.repo.GetPersonByID(id)

	if err != nil {
		Log.Info("Failed to get person by ID", "id", idStr, "error", err)

		return entities.Person{}, err
	}

	Log.Info("Person received successfully", "id", idStr)

	return person, nil
}

func (p *PersonService) GetPeopleByFilter(name, surname, ageStr, gender, nationality, page, limit string, patronymic *string) ([]entities.Person, error) {
	age, err := strconv.Atoi(ageStr)

	if err != nil {
		Log.Info("Invalid age format", "age", ageStr, "error", err)

		return nil, errors.New("invalid age format")
	}

	person := entities.Person{
		Name: name,
		Surname: surname,
		Patronymic: patronymic,
		Age: age,
		Gender: gender,
		Nationality: nationality,
	}

	people, err := p.repo.GetPeopleByFilter(person, page, limit)

	if err != nil {
		Log.Info("Failed to get people by age", "age", ageStr, "error", err)

		return nil, err
	}

	Log.Info("People received successfully by age", "age", ageStr, "page", page, "limit", limit)

	return people, nil
}

func (p *PersonService) DeletePersonByID(idStr string) error {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		Log.Info("Invalid ID format", "id", idStr, "error", err)
		return errors.New("invalid id format")
	}

	err = p.repo.DeletePersonByID(id)
	if err != nil {
		Log.Info("Failed to delete person", "id", id, "error", err)
		return err
	}

	Log.Info("Person deleted successfully", "id", id)

	return nil
}