package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/agl/fio/internal/domain/entities"
	. "github.com/agl/fio/pkg/logger"
)

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) UpdatePersonByID(id int, p entities.Person) error {
    tx, err := r.db.Begin()
    if err != nil {
        Log.Error("Failed to begin transaction", "error", err)

        return fmt.Errorf("failed to start transaction: %w", err)
    }
    defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

    var (
        currName, currSurname, currPatronymic sql.NullString
        currAge                           sql.NullInt64
        currGender, currNationality      string
        currGenderID, currNationalityID  int
    )
    query := `
        SELECT p.name, p.surname, p.patronymic, p.age,
               g.id, g.gender,
               n.id, n.nationality
        FROM people p
        JOIN genders g ON p.gender_id = g.id
        JOIN nationalities n ON p.nationality_id = n.id
        WHERE p.id = $1` 
    if err := tx.QueryRow(query, id).Scan(
        &currName, &currSurname, &currPatronymic, &currAge,
        &currGenderID, &currGender,
        &currNationalityID, &currNationality,
    ); err != nil {
        Log.Info("Failed to fetch current person data", "id", id, "error", err)
        return fmt.Errorf("failed to fetch current person data: %w", err)
    }

    newGenderID := currGenderID
    if p.Gender != "" && p.Gender != currGender {
        var gid int
        err = tx.QueryRow(`SELECT id FROM genders WHERE gender = $1`, p.Gender).Scan(&gid)
        if err == sql.ErrNoRows {
            err = tx.QueryRow(`INSERT INTO genders(gender) VALUES($1) RETURNING id`, p.Gender).Scan(&gid)
            if err != nil {
                Log.Info("Failed to insert new gender", "gender", p.Gender, "error", err)
                return fmt.Errorf("failed to insert gender: %w", err)
            }
        } else if err != nil {
            Log.Info("Failed to query gender existence", "gender", p.Gender, "error", err)
            return fmt.Errorf("failed to check gender: %w", err)
        }
        newGenderID = gid
    }

    newNationalityID := currNationalityID
    if p.Nationality != "" && p.Nationality != currNationality {
        var nid int
        err = tx.QueryRow(`SELECT id FROM nationalities WHERE nationality = $1`, p.Nationality).Scan(&nid)
        if err == sql.ErrNoRows {
            err = tx.QueryRow(`INSERT INTO nationalities(nationality) VALUES($1) RETURNING id`, p.Nationality).Scan(&nid)
            if err != nil {
                Log.Info("Failed to insert new nationality", "nationality", p.Nationality, "error", err)
                return fmt.Errorf("failed to insert nationality: %w", err)
            }
        } else if err != nil {
            Log.Info("Failed to query nationality existence", "nationality", p.Nationality, "error", err)
            return fmt.Errorf("failed to check nationality: %w", err)
        }
        newNationalityID = nid
    }

    setClauses := []string{}
    args := []any{}
    argPos := 1

    if p.Name != "" && p.Name != currName.String {
        setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
        args = append(args, p.Name)
        argPos++
    }
    if p.Surname != "" && p.Surname != currSurname.String {
        setClauses = append(setClauses, fmt.Sprintf("surname = $%d", argPos))
        args = append(args, p.Surname)
        argPos++
    }
    if p.Patronymic != nil && *p.Patronymic != currPatronymic.String {
        setClauses = append(setClauses, fmt.Sprintf("patronymic = $%d", argPos))
        args = append(args, *p.Patronymic)
        argPos++
    }
    if p.Age != 0 && int64(p.Age) != currAge.Int64 {
        setClauses = append(setClauses, fmt.Sprintf("age = $%d", argPos))
        args = append(args, p.Age)
        argPos++
    }
    if newGenderID != currGenderID {
        setClauses = append(setClauses, fmt.Sprintf("gender_id = $%d", argPos))
        args = append(args, newGenderID)
        argPos++
    }
    if newNationalityID != currNationalityID {
        setClauses = append(setClauses, fmt.Sprintf("nationality_id = $%d", argPos))
        args = append(args, newNationalityID)
        argPos++
    }

    if len(setClauses) > 0 {
        args = append(args, id)
        query = fmt.Sprintf("UPDATE people SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argPos)
        result, err := tx.Exec(query, args...)
        if err != nil {
            Log.Info("Failed to update people", "id", id, "error", err)
            return fmt.Errorf("failed to update person: %w", err)
        }
        rows, err := result.RowsAffected()
        if err != nil {
            Log.Info("Failed to retrieve affected rows", "id", id, "error", err)
            return fmt.Errorf("failed to retrieve affected rows: %w", err)
        }
        if rows == 0 {
            Log.Info("No person record updated", "id", id)
            return sql.ErrNoRows
        }
    }

    if err := tx.Commit(); err != nil {
        Log.Error("Failed to commit transaction", "error", err)
        return fmt.Errorf("transaction commit failed: %w", err)
    }

    Log.Info("Person and related data updated successfully", "id", id, "updates", p)
    return nil
}

func (r *PersonRepository) CreatePerson(person entities.Person) (int, error) {
	tx, err := r.db.Begin()
    if err != nil {
        Log.Error("Failed to begin transaction", "error", err)
		
        return 0, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var genderID int
	err = r.db.QueryRow(`SELECT id FROM genders WHERE gender = $1`, person.Gender).Scan(&genderID)
	if err != nil {
		err = r.db.QueryRow(`INSERT INTO genders (gender) VALUES ($1) RETURNING id`, person.Gender).Scan(&genderID)
		if err != nil {
			Log.Info("Failed to insert gender", "gender", person.Gender, "error", err)
			return 0, fmt.Errorf("failed to insert gender: %w", err)
		}
	}

	var nationalityID int
	err = r.db.QueryRow(`SELECT id FROM nationalities WHERE nationality = $1`, person.Nationality).Scan(&nationalityID)
	if err != nil {
		err = r.db.QueryRow(`INSERT INTO nationalities (nationality) VALUES ($1) RETURNING id`, person.Nationality).Scan(&nationalityID)
		if err != nil {
			Log.Info("Failed to insert nationality", "nationality", person.Nationality, "error", err)
			return 0, fmt.Errorf("failed to insert nationality: %w", err)
		}
	}

	query := `
		INSERT INTO people (name, surname, patronymic, age, gender_id, nationality_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err = r.db.QueryRow(query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		genderID,
		nationalityID,
	).Scan(&id)

	if err != nil {
		Log.Info("Failed to insert person", "person", person, "error", err)
		return 0, fmt.Errorf("failed to insert person: %w", err)
	}

	if err := tx.Commit(); err != nil {
        Log.Error("Failed to commit transaction", "error", err)

        return 0, fmt.Errorf("transaction commit failed: %w", err)
    }

	Log.Info("Person created successfully", "person", person, "id", id)
	return id, nil
}

func (r *PersonRepository) GetPersonByID(id int) (entities.Person, error) {
	query := `
		SELECT p.name, p.surname, p.patronymic, p.age, g.gender AS gender, n.nationality AS nationality
		FROM people p
		JOIN genders g ON p.gender_id = g.id
		JOIN nationalities n ON p.nationality_id = n.id
		WHERE p.id = $1
	`

	row := r.db.QueryRow(query, id)

	var p entities.Person
	err := row.Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)

	if err != nil {
		Log.Info("Failed to get person by ID", "id", id, "error", err)

		return entities.Person{}, err
	}

	Log.Info("Person retrieved successfully", "id", id, "person", p)

	return p, nil
}

func (r *PersonRepository) GetPeopleByFilter(filter entities.Person, page, limit string) ([]entities.Person, error) {
	query := `
		SELECT p.name, p.surname, p.patronymic, p.age, g.gender, n.nationality
		FROM people p
		JOIN genders g ON p.gender_id = g.id
		JOIN nationalities n ON p.nationality_id = n.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIdx := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND p.name = $%d", argIdx)
		args = append(args, filter.Name)
		argIdx++
	}
	if filter.Surname != "" {
		query += fmt.Sprintf(" AND p.surname = $%d", argIdx)
		args = append(args, filter.Surname)
		argIdx++
	}
	if filter.Patronymic != nil && *filter.Patronymic != "" {
		query += fmt.Sprintf(" AND p.patronymic = $%d", argIdx)
		args = append(args, *filter.Patronymic)
		argIdx++
	}
	if filter.Age != 0 {
		query += fmt.Sprintf(" AND p.age = $%d", argIdx)
		args = append(args, filter.Age)
		argIdx++
	}
	if filter.Gender != "" {
		query += fmt.Sprintf(" AND g.gender = $%d", argIdx)
		args = append(args, filter.Gender)
		argIdx++
	}
	if filter.Nationality != "" {
		query += fmt.Sprintf(" AND n.nationality = $%d", argIdx)
		args = append(args, filter.Nationality)
		argIdx++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, page)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		Log.Info("Failed to query people with filters", "filter", filter, "error", err)
		return nil, err
	}
	defer rows.Close()

	var people []entities.Person
	for rows.Next() {
		var p entities.Person
		err := rows.Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
		if err != nil {
			Log.Info("Failed to scan row", "error", err)
			return nil, err
		}
		people = append(people, p)
	}

	Log.Info("People retrieved with filters", "filter", filter, "count", len(people))
	return people, nil
}

func (r *PersonRepository) DeletePersonByID(id int) error {
	Log.Debug("Deleting person by ID", "ID", id)

	tx, err := r.db.Begin()
	if err != nil {
		Log.Error("Failed to begin transaction", "error", err)
        return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var genderID, nationalityID int

	err = r.db.QueryRow(`SELECT gender_id, nationality_id FROM people WHERE id = $1`, id).Scan(&genderID, &nationalityID)
	if err != nil {
		Log.Info("Failed to get gender/nationality ID for person", "ID", id, "error", err)

		return err
	}

	_, err = r.db.Exec(`DELETE FROM people WHERE id = $1`, id)
	if err != nil {
		Log.Info("Failed to delete person", "ID", id, "error", err)

		return err
	}

	_, _ = r.db.Exec(`
		DELETE FROM genders 
		WHERE id = $1
		`, genderID)

	_, _ = r.db.Exec(`
		DELETE FROM nationalities 
		WHERE id = $1
		`, nationalityID)

	if err := tx.Commit(); err != nil {
        Log.Error("Failed to commit transaction", "error", err)

        return fmt.Errorf("transaction commit failed: %w", err)
    }

	Log.Debug("Successfully deleted person and checked for unused gender/nationality", "ID", id)

	return nil
}