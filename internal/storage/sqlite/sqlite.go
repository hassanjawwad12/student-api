package sqlite

import (
	"database/sql"
	"github.com/hassanjawwad12/student-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	//the pattern is that err is always the second argument returned

	//database and the storage path is given to the function
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS (
	 id INTEGER PRIMARY KEY AUTOINCREMENT,
	 name TEXT ,
	 email TEXT ,
	 age INTEGER)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

// this method is now attched to the struct
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	//We prepare the data before binding it, this saves it from SQL injection
	stmt, err := s.Db.Prepare("INSERT INTO STUDENTS (name,email,age) VALUES (?,?,?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}
	return lastId, nil
}
