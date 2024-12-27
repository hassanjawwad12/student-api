package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hassanjawwad12/student-api/internal/storage"
	"github.com/hassanjawwad12/student-api/internal/types"
	"github.com/hassanjawwad12/student-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"
)

// the storage being passed here is the dependency injection which makes it extensive
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")
		var student types.Student

		//decode the json data inside the student variable
		err := json.NewDecoder(r.Body).Decode(&student)

		//eof is no more inputs are available (empty body error)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		//general error
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validate the request , follow a zero trust policy
		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})

		// Convert the string to a byte slice before writing
		//w.Write([]byte("Welcome to student API"))
	}
}
