package api_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
	"github.com/stretchr/testify/assert"
)

func TestShouldSaveUsersToDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	users := []db.User{ // Users that will be saved to DB
		{
			Name:        "John Doe",
			PhoneNumber: "18001234567",
			Country:     "US",
			City:        "New York City",
			ID:          1,
		},
		{
			Name:        "Florida Man",
			PhoneNumber: "18002234567",
			Country:     "US",
			City:        "Florida City",
			ID:          2,
		},
	}

	// Create a request that will save the users
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/users", dbUsersToCSV(users))
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	// Create a recorder to track what is happening in the router
	recorder := httptest.NewRecorder()

	// Prepare the server
	database := db.NewInMemoryDB()
	ginRouter := api.NewGinRouter(api.NewServer(database))

	// Send the request and check the DB
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, users, database.Users)
}

func TestShouldRejectCreateUsersWrongContentType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/users", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)
}

func TestShouldRejectCreateUsersBadCSV(t *testing.T) {
	t.Parallel()

	const badCSV = "notid,blablabla,,,"

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/users", strings.NewReader(badCSV))
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestShouldRejectCreateUsersNoBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/users", nil)
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestShouldRejectCreateUsersBadMethod(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func dbUsersToCSV(users []db.User) io.Reader {
	final := ""

	for _, user := range users {
		final += fmt.Sprintf("%d,%s,%s,%s,%s\n", user.ID, user.Name, user.PhoneNumber, user.Country, user.City)
	}

	return strings.NewReader(final)
}

func TestMain(m *testing.M) {
	logging.GlobalLogger = logging.StdLogger{}

	m.Run()
}
