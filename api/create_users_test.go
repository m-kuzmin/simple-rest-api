package api_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/testutil"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/stretchr/testify/assert"
)

func TestShouldSaveUsersToDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	users := []db.User{ // Users that will be saved to DB
		{
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
		{
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
	}

	// Create a request that will save the users
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", dbUsersToCSV(users))
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

func TestShouldRejectCreateUsersNoContentType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", nil)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", strings.NewReader(badCSV))
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestShouldRejectCreateUsersNilBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", nil)
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestShouldRejectCreateUsersEmptyBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users", strings.NewReader(""))
	assert.Nil(t, err)
	req.Header.Set("content-type", "text/csv")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestShouldCreateUsersWithFileUpload(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	users := []db.User{ // Users that will be saved to DB
		{
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
		{
			ID:          rand.Int63(),
			Name:        testutil.GenerateRandomAlphaOnlyString(10),
			PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
			Country:     testutil.GenerateRandomAlphaOnlyString(10),
			City:        testutil.GenerateRandomAlphaOnlyString(10),
		},
	}
	csvReader := dbUsersToCSV(users)

	buf := bytes.Buffer{}
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "file")
	assert.NoError(t, err)

	_, err = io.Copy(part, csvReader)
	assert.NoError(t, err)

	assert.NoError(t, writer.Close())

	t.Logf("buf = %s", buf.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users/upload", &buf)
	assert.Nil(t, err)

	req.Header.Set("content-type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	database := db.NewInMemoryDB()
	ginRouter := api.NewGinRouter(api.NewServer(database))

	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, users, database.Users)
}

func TestShouldRejectCreateUsersUploadNoContentType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users/upload", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)
}

func TestShouldRejectCreateUsersUploadBadCSV(t *testing.T) {
	t.Parallel()

	const badCSV = ",foo,,,bar,,,\n%$*()#@$"

	buf := bytes.Buffer{}
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "file")
	assert.NoError(t, err)

	_, err = part.Write([]byte(badCSV))
	assert.NoError(t, err)

	assert.NoError(t, writer.Close())

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users/upload", &buf)
	assert.Nil(t, err)

	req.Header.Set("content-type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestShouldRejectCreateUsersUploadNilBody(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users/upload", nil)
	assert.Nil(t, err)
	req.Header.Set("content-type", "multipart/form-data")

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestShouldRejectCreateUsersUploadEmptyFile(t *testing.T) {
	t.Parallel()

	buf := bytes.Buffer{}
	writer := multipart.NewWriter(&buf)

	_, err := writer.CreateFormFile("file", "file") // we never write to this file anything
	assert.NoError(t, err)

	assert.NoError(t, writer.Close())

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/users/upload", &buf)
	assert.Nil(t, err)

	req.Header.Set("content-type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	ginRouter := api.NewGinRouter(api.NewServer(db.NewInMemoryDB()))
	ginRouter.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func dbUsersToCSV(users []db.User) io.Reader {
	final := ""

	for _, user := range users {
		final += fmt.Sprintf("%d,%s,%s,%s,%s\n", user.ID, user.Name, user.PhoneNumber, user.Country, user.City)
	}

	return strings.NewReader(final)
}
