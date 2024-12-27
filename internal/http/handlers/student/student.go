package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/hassanjawwad12/student-api/internal/storage"
	"github.com/hassanjawwad12/student-api/internal/types"
	"github.com/hassanjawwad12/student-api/internal/utils/response"
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

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id") //this name must be same as u gave in the main.go file handlerfunc api calling
		slog.Info("getting the student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusOK, students)

	}
}

func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("deleting a student")

		idParam := r.PathValue("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid student ID"})
			return
		}

		// Delete the student from the database
		err = storage.DeleteStudent(id)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete student"})
			return
		}

		// Return a success response
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "student deleted successfully"})
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("updating a student")

		// Extract the student ID from the URL path using r.PathValue
		idParam := r.PathValue("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid student ID"})
			return
		}

		// Decode the request body into a map
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Fetch the existing student to get current values
		existingStudent, err := storage.GetStudentById(id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "student not found"})
			return
		}

		// Apply updates to the existing student
		name := existingStudent.Name
		if n, ok := updates["name"]; ok {
			name = n.(string)
		}

		email := existingStudent.Email
		if e, ok := updates["email"]; ok {
			email = e.(string)
		}

		age := existingStudent.Age
		if a, ok := updates["age"]; ok {
			age = int(a.(float64)) // JSON numbers are decoded as float64
		}

		// Update the student in the database
		err = storage.UpdateStudent(id, name, email, age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update student"})
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "student updated successfully"})
	}
}
