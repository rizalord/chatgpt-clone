package repository

import (
	"regexp"
	"testing"
	"time"
	"user-service/app/model/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var email = "test@example.com"

func TestUserRepositoryImpl_FindByEmail(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
        WithArgs(email, 1).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
            AddRow(1, "Test User", email, time.Now(), time.Now()))

    var user entity.User
    err := repo.FindByEmail(db, &user, email)

    assert.NoError(t, err)
    assert.Equal(t, email, user.Email)
}

func TestUserRepositoryImpl_FindByEmail_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
        WithArgs(email, 1).
        WillReturnError(gorm.ErrRecordNotFound)

    var user entity.User
    err := repo.FindByEmail(db, &user, email)

    assert.Error(t, err)
}

func TestUserRepositoryImpl_Create(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
        WithArgs("Test User", email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    user := &entity.User{
        Name:      "Test User",
        Email:     "test@example.com",
        Password:  func(s string) *string { return &s }("password"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, user)

    assert.NoError(t, err)
    assert.Equal(t, uint(1), user.ID)
}

func TestUserRepositoryImpl_Create_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
        WithArgs("Test User", email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnError(gorm.ErrInvalidData)

    user := &entity.User{
        Name:      "Test User",
        Email:     email,
        Password:  func(s string) *string { return &s }("password"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, user)

    assert.Error(t, err)
}

func TestUserRepositoryImpl_Update(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"password"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`)).
        WithArgs("Updated User", "updated@example.com", "password", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
        WillReturnResult(sqlmock.NewResult(1, 1))

    user := &entity.User{
        ID:        1,
        Name:      "Updated User",
        Email:     "updated@example.com",
        Password:  func(s string) *string { return &s }("password"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Update(db, user)

    assert.NoError(t, err)
}

func TestUserRepositoryImpl_Update_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepositoryImpl()

    mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"email"=$2,"password"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`)).
        WithArgs("Updated User", "updated@example.com", "password", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
        WillReturnError(gorm.ErrInvalidData)

    user := &entity.User{
        ID:        1,
        Name:      "Updated User",
        Email:     "updated@example.com",
        Password:  func(s string) *string { return &s }("password"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Update(db, user)

    assert.Error(t, err)
}

