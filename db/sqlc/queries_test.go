package sqlc_test

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"

	"github.com/docker/docker/testutil"
	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateAccount(t *testing.T) {
	t.Parallel()

	userID := rand.Int63()
	t.Log("user id:", userID)

	ctx := context.Background()
	arg := sqlc.CreateUserParams{
		ID:          userID,
		Name:        testutil.GenerateRandomAlphaOnlyString(10),
		PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
		Country:     testutil.GenerateRandomAlphaOnlyString(10),
		City:        testutil.GenerateRandomAlphaOnlyString(10),
	}

	err := testQueries.CreateUser(ctx, arg)
	assert.NoError(t, err, "While creating the user")

	err = testQueries.DeleteUserByID(ctx, userID)
	assert.NoError(t, err, "While deleting the user")
}

func TestShouldFindAtLeastOneRowWithEmptyQuery(t *testing.T) {
	t.Parallel()

	userID := rand.Int63()
	t.Log("user id:", userID)

	ctx := context.Background()
	createParams := sqlc.CreateUserParams{
		ID:          userID,
		Name:        testutil.GenerateRandomAlphaOnlyString(10),
		PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
		Country:     testutil.GenerateRandomAlphaOnlyString(10),
		City:        testutil.GenerateRandomAlphaOnlyString(10),
	}

	err := testQueries.CreateUser(ctx, createParams)
	assert.NoError(t, err, "While creating the user")

	t.Log("Created user")

	findParams := sqlc.SearchUsersParams{}

	users, err := testQueries.SearchUsers(ctx, findParams)
	assert.NoError(t, err, "While searching for the user")

	t.Log("Searched for user")

	// Because the tests are run in parallel we dont know if there will be more rows than the one we inserted.
	assert.Contains(t, users, sqlc.User(createParams))

	err = testQueries.DeleteUserByID(ctx, userID)
	assert.NoError(t, err, "While deleting the user")

	t.Log("Deleted user")
}

func TestShouldFindRowWithMissingParams(t *testing.T) {
	t.Parallel()

	userID := rand.Int63()
	t.Log("user id:", userID)

	ctx := context.Background()
	createParams := sqlc.CreateUserParams{
		ID:          userID,
		Name:        testutil.GenerateRandomAlphaOnlyString(10),
		PhoneNumber: testutil.GenerateRandomAlphaOnlyString(10),
		Country:     testutil.GenerateRandomAlphaOnlyString(10),
		City:        testutil.GenerateRandomAlphaOnlyString(10),
	}

	err := testQueries.CreateUser(ctx, createParams)
	assert.NoError(t, err, "While creating the user")

	t.Log("Created user")

	findParams := sqlc.SearchUsersParams{
		Column2: sql.NullString{Valid: true, String: createParams.PhoneNumber},
		Column3: sql.NullString{Valid: true, String: createParams.Country},
	}

	users, err := testQueries.SearchUsers(ctx, findParams)
	assert.NoError(t, err, "While searching for the user")

	t.Log("Searched for the user")

	assert.Equal(t, users, []sqlc.User{sqlc.User(createParams)})

	err = testQueries.DeleteUserByID(ctx, userID)
	assert.NoError(t, err, "While deleting the user")

	t.Log("Deleted user")
}
