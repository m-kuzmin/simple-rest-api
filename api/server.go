package api

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

// @Summary Add users to database
// @Description Add users to database by uploading a CSV file
// @Accept text/csv
// @Produce json
// @Success 201
// @Router /users [put]
func (s *Server) CreateOrUpdateUsers(ctx *gin.Context) {
	tape := logging.NewTape(
		logging.DebugLevel,
		logging.NewPrefixedLogger(logging.GlobalLogger, "(Tape (APICall PUT /users))"),
		logging.ErrorLevel,
		logging.NewPrefixedLogger(logging.GlobalLogger, "(APICall PUT /users)"),
	)

	tape.Debugf("%#v", ctx.Request)

	if ctx.ContentType() != "text/csv" {
		tape.Errorf("Wrong content type: %q", ctx.ContentType())
		errorResponse(ctx, http.StatusUnsupportedMediaType, `Expected Content-Type header to be "text/csv"`)

		return
	}

	if ctx.Request.Body == nil {
		tape.Errorf("Empty body")
		errorResponse(ctx, http.StatusUnprocessableEntity, "Empty CSV file not allowed")

		return
	}

	csvReader := csv.NewReader(ctx.Request.Body)

	users, err := ParseUsersCSV(csvReader)
	if err != nil {
		tape.Errorf("CSV parsing error: %s", err)
		errorResponsef(ctx, http.StatusUnprocessableEntity, "CSV parsing error: %s", err)

		return
	}

	if len(users) == 0 {
		tape.Errorf("Empty users list")
		errorResponse(ctx, http.StatusUnprocessableEntity, "User CSV file must contain at least one user")

		return
	}

	tape.Debugf("Users that will be added to DB: %v", users)

	err = s.db.CreateUsers(context.Background(), users)
	if err != nil {
		tape.Errorf("DB error while calling CreateUsers: %s", err)
		errorResponsef(ctx, http.StatusInternalServerError, "Database error: %s", err)

		return
	}

	tape.Infof("Returning StatusCreated")
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
