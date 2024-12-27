package storage

import (
	"github.com/hassanjawwad12/student-api/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteStudent(id int64) error
	UpdateStudent(id int64, name string, email string, age int) error
}
