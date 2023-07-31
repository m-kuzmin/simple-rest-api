package sqlc_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/m-kuzmin/simple-rest-api/db/sqlc"
)

func TestShouldCreateAccount(t *testing.T) {
	t.Parallel()

	userID := rand.Int63() //nolint:gosec // We only need this to prevent two tests from placing the same ID in the DB.
	t.Log("user id:", userID)

	ctx := context.Background()
	arg := sqlc.CreateUserParams{
		ID:          userID,
		Name:        "John Doe",
		PhoneNumber: "18001234567",
		Country:     "US",
		City:        "New York",
	}

	err := testQueries.CreateUser(ctx, arg)
	if err != nil {
		t.Logf("While creating the user: %s", err)
		t.Fail()
	}

	err = testQueries.DeleteUserByID(ctx, userID)
	if err != nil {
		t.Logf("While deleting the user: %s", err)
		t.Fail()
	}
}
