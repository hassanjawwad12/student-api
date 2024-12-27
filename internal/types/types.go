package types

type Student struct {
	Id    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"required"`
	Email string `json:"email" validate:"required"`
}
