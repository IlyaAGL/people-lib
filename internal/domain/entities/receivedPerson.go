package entities

// ReceivedPerson represents a received person entity
// @Description Person information sent by client
type ReceivedPerson struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
}