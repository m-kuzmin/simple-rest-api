package api

import (
	"context"
	"encoding/csv"
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
// .
func (s *Server) CreateOrUpdateUsers(ctx *gin.Context) {
	if ctx.ContentType() != "text/csv" {
		logging.Errorf("(APICall PUT /users) Wrong content type: %q", ctx.ContentType())
		errorResponse(ctx, http.StatusUnsupportedMediaType, `Expected Content-Type header to be "text/csv"`)

		return
	}

	if ctx.Request.Body == nil {
		logging.Errorf("(APICall PUT /users) Empty body")
		errorResponse(ctx, http.StatusUnprocessableEntity, "Empty CSV file not allowed")

		return
	}

	records, err := csv.NewReader(ctx.Request.Body).ReadAll()
	if err != nil {
		logging.Errorf("(APICall PUT /users) CSV parsing error: %s", err)
		errorResponsef(ctx, http.StatusUnprocessableEntity, "Could not parse CSV: %s", err.Error())

		return
	}

	users := make([]db.User, len(records))

	for i, rec := range records {
		id, err := strconv.ParseInt(rec[0], 10, 0) //nolint:govet // Shadowing err is ok
		if err != nil {
			logging.Errorf("(APICall PUT /users) CSV entry %d: User.ID isnt a number: %s", i, err)
			errorResponsef(ctx, http.StatusUnprocessableEntity, "CSV entry %d: ID is not a number: %q", i, rec[0])

			return
		}

		users[i] = db.User{
			Name:        rec[1],
			PhoneNumber: rec[2],
			Country:     rec[3],
			City:        rec[4],
			ID:          id,
		}
	}

	logging.Debugf("Users that will be added to DB: %v", users)

	err = s.db.CreateUsers(context.Background(), users)
	if err != nil {
		logging.Errorf("(APICall PUT /users) DB error while calling CreateUsers: %s", err)
		errorResponsef(ctx, http.StatusInternalServerError, "Failed to save users to database: %s", err)

		return
	}

	logging.Infof("(APICall PUT /users) Returning StatusCreated")
	okResponse(ctx, http.StatusCreated)
}
