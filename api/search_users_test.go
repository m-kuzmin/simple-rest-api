package api_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/docker/docker/testutil"
	"github.com/gin-gonic/gin"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/stretchr/testify/assert"
)

func TestShouldSearchUsersByOneColumn(t *testing.T) {
	t.Parallel()

	database := db.NewInMemoryDB()
	ctx := context.Background()
	database.Users = []db.User{ // Put a user into the DB
		{
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		}, {
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
	}

	// Prepare a request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users/search", nil)
	assert.Nil(t, err)

	query := url.Values{}
	query.Set("name", database.Users[0].Name)
	req.URL.RawQuery = query.Encode()

	// Create a recorder to check the request
	recorder := httptest.NewRecorder()

	// Prepare the server
	server := api.NewServer(database)
	router := api.NewGinRouter(server)

	// Send the request and check the results
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Contains(t, recorder.Header().Get("content-type"), "application/json")

	jsonBody, err := json.Marshal(gin.H{"ok": true, "results": []db.User{database.Users[0]}})
	assert.NoError(t, err)

	assert.Equal(t, string(jsonBody), recorder.Body.String())
}

func TestShouldSearchUsersByAllColumns(t *testing.T) {
	t.Parallel()

	database := db.NewInMemoryDB()
	ctx := context.Background()
	sameName := testutil.GenerateRandomAlphaOnlyString(10)
	samePhoneNumber := testutil.GenerateRandomAlphaOnlyString(10)
	database.Users = []db.User{ // Put a user into the DB
		{
			ID:          rand.Int63(),
			Name:        sameName,
			PhoneNumber: samePhoneNumber,
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		}, {
			ID:          rand.Int63(),
			Name:        sameName,
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		}, {
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: samePhoneNumber,
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
	}

	// Prepare a request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users/search", nil)
	assert.Nil(t, err)

	query := url.Values{}
	query.Set("name", database.Users[0].Name)
	query.Set("phone_number", database.Users[0].PhoneNumber)
	query.Set("country", database.Users[0].Country)
	query.Set("city", database.Users[0].City)
	req.URL.RawQuery = query.Encode()

	// Create a recorder to check the request
	recorder := httptest.NewRecorder()

	// Prepare the server
	server := api.NewServer(database)
	router := api.NewGinRouter(server)

	// Send the request and check the results
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusFound, recorder.Code)
	assert.Contains(t, recorder.Header().Get("content-type"), "application/json")

	jsonBody, err := json.Marshal(gin.H{"ok": true, "results": []db.User{database.Users[0]}})
	assert.NoError(t, err)

	assert.Equal(t, string(jsonBody), recorder.Body.String())
}

func TestEmptrySearchQuery(t *testing.T) {
	t.Parallel()

	database := db.NewInMemoryDB()
	ctx := context.Background()

	// Prepare a request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/users/search", nil)
	assert.Nil(t, err)

	// Create a recorder to check the request
	recorder := httptest.NewRecorder()

	// Prepare the server
	server := api.NewServer(database)
	router := api.NewGinRouter(server)

	// Send the request and check the results
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Header().Get("content-type"), "application/json")

	jsonBody, err := json.Marshal(gin.H{"ok": false, "error": "Empty search criteria"})
	assert.NoError(t, err)

	assert.Equal(t, string(jsonBody), recorder.Body.String())
}
