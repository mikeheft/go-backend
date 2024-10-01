package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mikeheft/go-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user, user2)
	require.WithinDuration(t, user.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, user2.UpdatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)
	fullname := util.RandomOwner()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: fullname,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, fullname, updatedUser.FullName)
}
