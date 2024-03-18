package actor

type Actor struct {
	Name      string `json:"name" notempty:"true"`
	Gender    string `json:"gender" validate:"oneof=man woman"`
	BirthDate string `json:"birth_date" notempty:"true"`
}
