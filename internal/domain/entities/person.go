package entities

// Person represents a person entity
// @Description Person information with age, gender and nationality
type Person struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         int     `json:"age"`
	Gender      string  `json:"gender"`
	Nationality string  `json:"nationality"`
}

type Nationalities struct {
	Countries []Nationality `json:"country"`
}

type Nationality struct {
	CountryID string `json:"country_id"`
}
