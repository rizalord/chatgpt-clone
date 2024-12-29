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

func TestMessageRepositoryImpl_FindAllByChatID(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewMessageRepositoryImpl()

    chatID := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE chat_id = $1`)).
        WithArgs(chatID).
        WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "content", "created_at", "updated_at"}).
            AddRow(1, chatID, "Hello", time.Now(), time.Now()))

    var messages []entity.Message
    err := repo.FindAllByChatID(db, &messages, chatID)

    assert.NoError(t, err)
    assert.Len(t, messages, 1)
    assert.Equal(t, chatID, messages[0].ChatID)
}

func TestMessageRepositoryImpl_FindAllByChatID_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewMessageRepositoryImpl()

    chatID := 1
    mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "messages" WHERE chat_id = $1`)).
        WithArgs(chatID).
        WillReturnError(gorm.ErrRecordNotFound)

    var messages []entity.Message
    err := repo.FindAllByChatID(db, &messages, chatID)

    assert.Error(t, err)
}

func TestMessageRepositoryImpl_Create(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewMessageRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "messages" ("chat_id","content","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
        WithArgs(1, "Hello", sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    message := &entity.Message{
        ChatID:    1,
        Content:   "Hello",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, message)

    assert.NoError(t, err)
    assert.Equal(t, uint(1), message.ID)
}

func TestMessageRepositoryImpl_Create_Error(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewMessageRepositoryImpl()

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "messages" ("chat_id","content","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
        WithArgs(1, "Hello", sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnError(gorm.ErrInvalidData)

    message := &entity.Message{
        ChatID:    1,
        Content:   "Hello",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := repo.Create(db, message)

    assert.Error(t, err)
}