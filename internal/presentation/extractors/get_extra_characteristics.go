package extractors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/agl/fio/internal/domain/entities"
	. "github.com/agl/fio/pkg/logger"
)

const timeout = 3 * time.Second

func GetExtraUserInfoByName(name string, received_person entities.ReceivedPerson) (entities.Person, error) {
	person := entities.Person{
		Name: received_person.Name,
		Surname: received_person.Surname,
		Patronymic: received_person.Patronymic,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := getAge(ctx, name, &person)
	if err != nil {
		Log.Info("Failed to get age", "name", name, "error", err)

		return entities.Person{}, fmt.Errorf("failed to get age: %w", err)
	}

	err = getGender(ctx, name, &person)
	if err != nil {
		Log.Info("Failed to get gender", "name", name, "error", err)

		return entities.Person{}, fmt.Errorf("failed to get gender: %w", err)
	}

	err = getNationality(ctx, name, &person)
	if err != nil {
		Log.Info("Failed to get nationality", "name", name, "error", err)

		return entities.Person{}, fmt.Errorf("failed to get nationality: %w", err)
	}

	return person, nil
}

func getAge(ctx context.Context, name string, person *entities.Person) error {
	agifu_url := os.Getenv("AGIFY_URL")

	url := agifu_url + "=" + name

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		Log.Info("Failed to create age request", "name", name, "error", err)

		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Info("Failed to get age from API", "name", name, "error", err)

		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Info("Failed to read response body for age", "name", name, "error", err)

		return err
	}

	if err := json.Unmarshal(body, person); err != nil {
		Log.Info("Failed to unmarshal age response", "name", name, "error", err)

		return err
	}

	Log.Info("Successfully retrieved person age", "age", person.Age)

	return nil
}

func getGender(ctx context.Context, name string, person *entities.Person) error {
	genderize_url := os.Getenv("GENDERIZE_URL")

	url := genderize_url + "=" + name

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		Log.Info("Failed to create gender request", "name", name, "error", err)

		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Info("Failed to get gender from API", "name", name, "error", err)

		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Info("Failed to read response body for gender", "name", name, "error", err)

		return err
	}

	if err := json.Unmarshal(body, person); err != nil {
		Log.Info("Failed to unmarshal gender response", "name", name, "error", err)

		return err
	}

	Log.Info("Successfully retrieved person gender", "gender", person.Gender)

	return nil
}

func getNationality(ctx context.Context, name string, person *entities.Person) error {
	nationalize_url := os.Getenv("NATIONALIZE_URL")

	url := nationalize_url + "=" + name

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		Log.Info("Failed to create nationality request", "name", name, "error", err)

		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Info("Failed to get nationality from API", "name", name, "error", err)

		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Info("Failed to read response body for nationality", "name", name, "error", err)

		return err
	}

	var result entities.Nationalities
	if err := json.Unmarshal(body, &result); err != nil {
		Log.Info("Failed to unmarshal nationality response", "name", name, "error", err)

		return err
	}

	if len(result.Countries) > 0 {
		person.Nationality = result.Countries[0].CountryID

		return nil
	}

	Log.Info("No nationality found", "name", name)

	return fmt.Errorf("no nationality found")
}
