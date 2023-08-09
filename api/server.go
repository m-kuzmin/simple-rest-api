package api

// @title Simple REST API
// @version 0.1.0
// @description A simple REST API server with PostgreSQL database

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
)

type Server struct {
	db db.Querier
}

func NewServer(db db.Querier) *Server {
	return &Server{db: db}
}

// @Tags Add user data to database
// @Router /users [post]
// @Summary Add users to database (Request body)
// @Description Add users to database by sending a CSV body. The ID must not be already in the database.
// @Accept text/csv
// @Param CSV body string true "CSV string"
// @Produce json
// @Success 201 {object} api.CreateUsersBody.responseOk "Data was saved to database."
// @Failure 400 {object} api.CreateUsersBody.responseErr "Empty body not allowed"
// @Failure 415 {object} api.CreateUsersBody.responseErr "Content-Type must be text/csv."
// @Failure 422 {object} api.CreateUsersBody.responseErr "CSV body has syntax errors"
// @Failure 500 {object} api.CreateUsersBody.responseErr "Errors from PostgreSQL"
func (s *Server) CreateUsersBody(ctx *gin.Context) {
	type responseOk struct {
		Ok bool `json:"ok" example:"true"`
	}

	type responseErr struct {
		Error string `json:"error" example:"Human readable error"`
		Ok    bool   `json:"ok" example:"false"`
	}

	tape := logging.NewTape(
		logging.DebugLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(Tape (APICall PUT /users))"),
		logging.ErrorLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(APICall PUT /users)"),
	)

	tape.Debugf("%v", ctx.Request)

	if ctx.ContentType() != "text/csv" {
		tape.Errorf("Wrong content type: %q", ctx.ContentType())
		errorResponse(ctx, http.StatusUnsupportedMediaType, "Content-Type must be text/csv")

		return
	}

	if ctx.Request.Body == nil {
		tape.Errorf("Empty body (request.body == nil)")
		errorResponse(ctx, http.StatusBadRequest, "Empty body not allowed")

		return
	}

	csvReader := csv.NewReader(ctx.Request.Body)

	users, err := ParseUsersCSV(csvReader)
	if err != nil {
		tape.Errorf("CSV parsing error: %s", err)
		errorResponsef(ctx, http.StatusUnprocessableEntity, "CSV parsing error: %s", err)

		return
	}

	s.saveUsersToDB(ctx, tape, users)
}

// @Tags Add user data to database
// @Router /users/upload [post]
// @Summary Add users to database (File upload)
// @Description Add users to database by uploading a CSV file. The ID must not be already in the database.
// @Accept multipart/form-data
// @Param file formData file true "CSV string"
// @Produce json
// @Success 201 {object} api.CreateUsersBody.responseOk "Data was saved to database"
// @Failure 400 {object} api.CreateUsersBody.responseErr "Empty body not allowed"
// @Failure 415 {object} api.CreateUsersBody.responseErr "Content-Type must be multipart/form-data"
// @Failure 422 {object} api.CreateUsersBody.responseErr "CSV body has syntax errors"
// @Failure 500 {object} api.CreateUsersBody.responseErr "Errors from PostgreSQL"
func (s *Server) CreateUsersUpload(ctx *gin.Context) {
	tape := logging.NewTape(
		logging.DebugLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(Tape (APICall PUT /users/upload))"),
		logging.ErrorLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(APICall PUT /users/upload)"),
	)

	tape.Debugf("%+v", ctx.Request)

	if ctx.ContentType() != "multipart/form-data" {
		tape.Errorf("Wrong content type: %q", ctx.ContentType())
		errorResponse(ctx, http.StatusUnsupportedMediaType, "Content-Type must be multipart/form-data")

		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		tape.Errorf("Failed to get form file: %s", err)
		errorResponse(ctx, http.StatusBadRequest,
			"File processing error: Attached file with name 'file' not found")

		return
	}

	file, err := fileHeader.Open()
	defer func() {
		if err = file.Close(); err != nil {
			tape.Errorf("Failed to close file")
		}
	}()

	if err != nil {
		tape.Errorf("Failed to open form file: %s", err)
		errorResponse(ctx, http.StatusInternalServerError, "File processing error: Failed to open attached file")

		return
	}

	csvReader := csv.NewReader(file)

	users, err := ParseUsersCSV(csvReader)
	if err != nil {
		tape.Errorf("CSV parsing error: %s", err)
		errorResponsef(ctx, http.StatusUnprocessableEntity, "CSV parsing error: %s", err)

		return
	}

	s.saveUsersToDB(ctx, tape, users)
}

func (s *Server) saveUsersToDB(ctx *gin.Context, logger logging.Logger, users []db.User) {
	if len(users) == 0 {
		logger.Errorf("Empty users list")
		errorResponse(ctx, http.StatusBadRequest, "User CSV file must contain at least one user")

		return
	}

	logger.Debugf("Users that will be added to DB: %v", users)

	err := s.db.CreateUsers(context.Background(), users)
	if err != nil {
		logger.Errorf("DB error while calling CreateUsers: %s", err)
		errorResponsef(ctx, http.StatusInternalServerError, "Database error: %s", err)

		return
	}

	logger.Infof("Returning StatusCreated")
	okResponse(ctx, http.StatusCreated)
}

/*
ParseUsersCSV parses the CSV file into a User list. If the CSV file has syntax errors returns (nil, err). If there is a
parsing error for one of the fields, returns all users parsed before the bad one and parsing error.
*/
func ParseUsersCSV(reader *csv.Reader) ([]db.User, error) {
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV records: %w", err)
	}

	users := make([]db.User, len(records))

	for i, rec := range records {
		id, err := strconv.ParseInt(rec[0], 10, 0)
		if err != nil {
			return users, fmt.Errorf("record %d: ID is not a number: %w", i, err)
		}

		users[i] = db.User{
			Name:        rec[1],
			PhoneNumber: rec[2],
			Country:     rec[3],
			City:        rec[4],
			ID:          id,
		}
	}

	return users, nil
}

// @Tags Search
// @Router /users/search [get]
// @Summary Search the database by column(s)
// @Description Returns all rows that contain the substring from the query in the respective column
// @Param name query string false "Reject all rows where name doesnt contain the substring"
// @Param phone_number query string false "Reject all rows where phone_number doesnt contain the substring"
// @Param country query string false "Reject all rows where country doesnt contain the substring"
// @Param city query string false "Reject all rows where city doesnt contain the substring"
// @Produce json
// @Success 302 {object} api.SearchUsers.responseOk "0+ JSON encoded objects"
// @Failure 400 {object} api.SearchUsers.responseErr "0 `?field=` criteria provided. Should have at least one."
// @Failure 500 {object} api.SearchUsers.responseErr "Errors from PostgreSQL"
func (s *Server) SearchUsers(ctx *gin.Context) {
	type responseOk struct {
		Ok      bool      `json:"ok" example:"true"`
		Results []db.User `json:"results"`
	}

	type responseErr struct {
		Error string `json:"error" example:"Human readable error"`
		Ok    bool   `json:"ok" example:"false"`
	}

	tape := logging.NewTape(
		logging.DebugLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(Tape (APICall GET /users/search))"),
		logging.ErrorLevel, logging.NewPrefixedLogger(logging.GlobalLogger, "(APICall GET /users/search)"),
	)

	tape.Debugf("Request: %v", ctx.Request)

	name := ctx.Query("name")
	phoneNumber := ctx.Query("phone_number")
	country := ctx.Query("country")
	city := ctx.Query("city")

	if name == "" && phoneNumber == "" && country == "" && city == "" {
		tape.Errorf("Empty search query")
		errorResponse(ctx, http.StatusBadRequest, "Empty search criteria")

		return
	}

	users, err := s.db.SearchUsers(context.Background(), name, phoneNumber, country, city)
	if err != nil {
		tape.Errorf("Database error: %s", err)
		errorResponsef(ctx, http.StatusInternalServerError, "Database error: %s", err)

		return
	}

	ctx.JSON(http.StatusFound, gin.H{"ok": true, "results": users})
}
