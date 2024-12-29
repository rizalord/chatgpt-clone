package usecase

import (
	"auth-service/app/helper"
	"auth-service/app/model/dto"
	"auth-service/app/model/entity"
	"auth-service/app/model/proto"
	"auth-service/app/repository"
	"context"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MessageUseCase interface {
	GetMessages(ctx context.Context, gmr *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error)
	CreateMessage(cmr proto.ChatService_CreateMessageServer) error
}

type MessageUseCaseImpl struct {
	DB              	*gorm.DB
	MessageRepository  	repository.MessageRepository
	ChatRepository     	repository.ChatRepository
	GenAI				*genai.Client
	Log             	*logrus.Logger
	Validate        	*validator.Validate
}

func NewMessageUseCaseImpl(
	db *gorm.DB,
	messageRepository repository.MessageRepository,
	chatRepository repository.ChatRepository,
	genai *genai.Client,
	log *logrus.Logger,
	validate *validator.Validate,
) *MessageUseCaseImpl {
	return &MessageUseCaseImpl{
		DB:             db,
		MessageRepository: messageRepository,
		ChatRepository: chatRepository,
		GenAI:			genai,
		Log:            log,
		Validate:       validate,
	}
}

func (c *MessageUseCaseImpl) GetMessages(ctx context.Context, gmr *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	// create request
	req := &dto.GetMessagesRequest{
		ChatID: uint(gmr.GetChatId()),
		UserID: uint(gmr.GetUserId()),
	}

	// validate request
	if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
		err := status.Newf(codes.InvalidArgument, "Validation errors: %v", errors)

		err, wde := err.WithDetails(gmr)
		if wde != nil {
			c.Log.Error("Failed to add details to error", wde)
			return nil, status.Error(codes.Internal, "An error occurred while adding details to error")
		}

		return nil, err.Err()
	}

	db := c.DB.WithContext(ctx)

	// get messages
	messages := new([]entity.Message)
	if err := c.MessageRepository.FindAllByChatID(db, messages, int(req.ChatID)); err != nil {
		c.Log.Error("Failed to get messages", err)
		return nil, status.Error(codes.NotFound, "Messages not found")
	}

	// create response
	res := &proto.GetMessagesResponse{
		Messages: make([]*proto.FullMessage, 0),
	}

	for _, message := range *messages {
		res.Messages = append(res.Messages, &proto.FullMessage{
			Id:     int64(message.ID),
			ChatId: int64(message.ChatID),
			UserId: int64(message.UserID),
			Role:  message.Role,
			Content: message.Content,
		})
	}

	return res, nil
}

func (c *MessageUseCaseImpl) CreateMessage(cmr proto.ChatService_CreateMessageServer) error {
	for {
		// receive message
		message, err := cmr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.Log.Error("Failed to receive client message", err)
			return status.Error(codes.Internal, "An error occurred while receiving message")
		}

		// create request
		req := &dto.CreateMessageRequest{
			ChatID: uint(message.GetChatId()),
			UserID: uint(message.GetUserId()),
			Message: message.GetMessage(),
		}

		// validate request
		if errors := helper.Validate(c.Validate, req); len(errors) > 0 {
			err := status.Newf(codes.InvalidArgument, "Validation errors: %v", errors)

			err, wde := err.WithDetails(message)
			if wde != nil {
				c.Log.Error("Failed to add details to error", wde)
				return status.Error(codes.Internal, "An error occurred while adding details to error")
			}

			return err.Err()
		}

		db := c.DB.WithContext(cmr.Context())
		tx := db.Begin()
		defer tx.Rollback()

		// Check if chat exists
		chat := new(entity.Chat)
		if err := c.ChatRepository.FindByID(tx, chat, int(req.ChatID)); err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create chat if not exists
				chat.UserID = req.UserID
				chat.Topic = nil

				if err := c.ChatRepository.Create(tx, chat); err != nil {
					c.Log.Error("Failed to create chat", err)
					return status.Error(codes.Internal, "An error occurred while creating chat")
				}
			} else {
				c.Log.Error("Failed to get chat", err)
				return status.Error(codes.Internal, "An error occurred while getting chat")
			}

		}

		// Initialize messages with history if exists
		dbMessages := new([]entity.Message)
		if err := c.MessageRepository.FindAllByChatID(db, dbMessages, int(req.ChatID)); err != nil {
			c.Log.Error("Failed to get history messages", err)
			return status.Error(codes.Internal, "An error occurred while getting messages")
		}

		messages := []genai.Content{
			{
				Parts: []genai.Part{
					genai.Text("You are a helpful chatbot. Answer user's chat with markdown format."),
				},
				Role: "user",
			},
		}

		for _, dbMessage := range *dbMessages {
			messages = append(messages, genai.Content{
				Parts: []genai.Part{
					genai.Text(dbMessage.Content),
				},
				Role: dbMessage.Role,
			})
		}

		// Start stream the messages
		model := c.GenAI.GenerativeModel("gemini-1.5-flash")
		cs := model.StartChat()

		cs.History = make([]*genai.Content, len(messages))
		for i := range messages {
			cs.History[i] = &messages[i]
		}

		cmr.Send(&proto.Part{
			ChatId: int64(chat.ID),
			Part:  "Starting chat...",
			Status: proto.Status_START,
		})

		iter := cs.SendMessageStream(cmr.Context(), genai.Text(req.Message))
		respMessage := ""
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.Log.Error("Failed to get next message", err)
				return status.Error(codes.Internal, "An error occurred while sending message")
			}

			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						respMessage += fmt.Sprintf("%v", part)
						cmr.Send(&proto.Part{
							ChatId: int64(req.ChatID),
							Part:  fmt.Sprintf("%v", part),
							Status: proto.Status_PROGRESS,
						})
					}
				}
			}
		}

		cmr.Send(&proto.Part{
			ChatId: int64(chat.ID),
			Part:  "Chat ended.",
			Status: proto.Status_END,
		})

		// create message
		userMsg := &entity.Message{
			ChatID: uint(chat.ID),
			UserID: uint(chat.UserID),
			Role:  "user",
			Content: req.Message,
		}

		respMsg := &entity.Message{
			ChatID: uint(chat.ID),
			UserID: uint(chat.UserID),
			Role:  "model",
			Content: respMessage,
		}

		if err := c.MessageRepository.Create(tx, userMsg); err != nil {
			c.Log.Error("Failed to create user message", err)
			return status.Error(codes.Internal, "An error occurred while creating user message")
		}

		if err := c.MessageRepository.Create(tx, respMsg); err != nil {
			c.Log.Error("Failed to create model message", err)
			return status.Error(codes.Internal, "An error occurred while creating model message")
		}

		// If length of history is lower than or equal to 3, create summarize topic
		if len(cs.History) <= 3 {
			summarizedTopic, err := c.SummarizeTopic(cmr.Context(), cs.History)
			if err != nil {
				c.Log.Error("Failed to get summarized topic", err)
				return status.Error(codes.Internal, "An error occurred while summarizing topic")
			}
	
			// Update chat with summarized topic
			chat.Topic = &summarizedTopic
			if err := c.ChatRepository.Update(tx, chat); err != nil {
				c.Log.Error("Failed to update chat", err)
				return status.Error(codes.Internal, "An error occurred while updating chat")
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return status.Error(codes.Internal, "An error occurred while committing transaction")
		}
	}

	return nil
}

func (c *MessageUseCaseImpl) SummarizeTopic(ctx context.Context, history []*genai.Content) (string, error) {
	model := c.GenAI.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	cs.History = make([]*genai.Content, len(history))
	copy(cs.History, history)

	res, err := cs.SendMessage(ctx, genai.Text("Get topic that user is asking, maximum 5 words, use user's language"))
	if err != nil {
		return "", err
	}

	summarizedTopic := ""
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				summarizedTopic += fmt.Sprintf("%v", part)
			}
		}
	}

	return summarizedTopic, nil
}