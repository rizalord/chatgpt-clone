package repository

import (
	"auth-service/app/model/entity"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestChatRepositoryImpl_FindByID(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    id := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chats" WHERE id = $1`)).
        WithArgs(id).
        WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "message", "created_at", "updated_at"}).
            AddRow(id, 1, "Hello", time.Now(), time.Now()))

    chat := &entity.Chat{}
    err := repo.FindByID(db, chat, id)

    assert.NoError(t, err)
    assert.Equal(t, id, chat.ID)
}

func TestChatRepositoryImpl_FindByID_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    id := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chats" WHERE id = $1`)).
        WithArgs(id).
        WillReturnError(gorm.ErrRecordNotFound)

    chat := &entity.Chat{}
    err := repo.FindByID(db, chat, id)

    assert.Error(t, err)
}

func TestChatRepositoryImpl_FindAllByUserID(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    userID := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chats" WHERE user_id = $1`)).
        WithArgs(userID).
        WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "message", "created_at", "updated_at"}).
            AddRow(1, userID, "Hello", time.Now(), time.Now()))

    var chats []entity.Chat
    err := repo.FindAllByUserID(db, &chats, userID)

    assert.NoError(t, err)
    assert.Len(t, chats, 1)
    assert.Equal(t, userID, chats[0].UserID)
}

func TestChatRepositoryImpl_FindAllByUserID_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    userID := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chats" WHERE user_id = $1`)).
        WithArgs(userID).
        WillReturnError(gorm.ErrRecordNotFound)

    var chats []entity.Chat
    err := repo.FindAllByUserID(db, &chats, userID)

    assert.Error(t, err)
}

func TestChatRepositoryImpl_Create(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "chats" ("user_id","message","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
        WithArgs(1, "Hello", sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    chat := &entity.Chat{
        UserID:    1,
        Topic:    func(s string) *string { return &s }("Hello"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, chat)

    assert.NoError(t, err)
    assert.Equal(t, uint(1), chat.ID)
}

func TestChatRepositoryImpl_Create_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "chats" ("user_id","message","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
        WithArgs(1, "Hello", sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnError(gorm.ErrInvalidData)

    chat := &entity.Chat{
        UserID:    1,
        Topic:    func(s string) *string { return &s }("Hello"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, chat)

    assert.Error(t, err)
}

func TestChatRepositoryImpl_Update(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    chat := &entity.Chat{
        ID:       1,
        UserID:    1,
        Topic:    func(s string) *string { return &s }("Hello"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chats" SET "user_id"=$1,"message"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
        WithArgs(chat.UserID, chat.Topic, chat.CreatedAt, chat.UpdatedAt, chat.ID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err := repo.Update(db, chat)

    assert.NoError(t, err)    
}

func TestChatRepositoryImpl_Update_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewChatRepositoryImpl()

    chat := &entity.Chat{
        ID:       1,
        UserID:    1,
        Topic:    func(s string) *string { return &s }("Hello"),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chats" SET "user_id"=$1,"message"=$2,"created_at"=$3,"updated_at"=$4 WHERE "id" = $5`)).
        WithArgs(chat.UserID, chat.Topic, chat.CreatedAt, chat.UpdatedAt, chat.ID).
        WillReturnError(gorm.ErrInvalidData)

    err := repo.Update(db, chat)

    assert.Error(t, err)
}