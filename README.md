# ChatGPT Clone with Gemini AI

A chatgpt clone web application that leverages Google's Gemini AI 1.5 Flash model to provide intelligent conversational responses. Built with microservices architecture using Go, Next.js, and gRPC.

## Demo

![ChatGPT Clone](https://i.imgur.com/9hyNmw8.gif)

You can access the live demo of the application [here](https://chatbot-web-service.vercel.app/).

## Architecture

The application is built using microservices architecture with the following services:

- **Web Service**: Next.js frontend application
- **API Gateway**: Go service that handles routing and communication between services
- **Auth Service**: Handles user authentication and authorization
- **Chat Service**: Manages chat functionality and Gemini AI integration
- **User Service**: Handles user management and profiles

Below is the diagram of the application architecture:
![ChatGPT Clone Architecture](https://i.imgur.com/kTx5TOg.png)

## Features

- Real-time chat functionality using Socket.IO
- Google OAuth authentication
- Conversational AI powered by Gemini 1.5 Flash
- Responsive web interface
- Chat history

## Prerequisites

- Docker and Docker Compose
- Go 1.23.2 or later
- Node.js 20 or later
- PostgreSQL
- Google Cloud Platform account for Gemini AI API

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/rizalord/chatgpt-clone.git
cd chatgpt-clone
```

2. Configure environment variables:
- Copy `.env.example` to `.env` in each service directory
- Set up required environment variables, including:
  - Database credentials
  - Gemini API key
  - Google OAuth credentials
  - JWT secrets

## Running the Application

### Run The Migration
Edit the database URL in the `migrate` command to match your database configuration. You will need the `migrate` package from [golang-migrate](https://github.com/golang-migrate/migrate) to run this command.
```bash
migrate -database "postgres://postgres@localhost:5432/chatgpt_clone?sslmode=disable" -path ./migrations up
```

### Using Docker Compose

```bash
docker compose up -d
```

### Using Local Development

To run the application locally without Docker, follow these steps:

```bash
# Start the gRPC services
cd services/user-service && go run cmd/worker/main.go
cd services/auth-service && go run cmd/worker/main.go
cd services/chat-service && go run cmd/worker/main.go

# Start the API Gateway
cd services/api-gateway && go run cmd/web/main.go

# Start the web frontend
cd services/web-service && npm install && npm run dev
```

## Project Structure
```
├── services/
│   ├── api-gateway/      # API Gateway service
│   ├── auth-service/     # Authentication service
│   ├── chat-service/     # Chat and Gemini AI service
│   ├── user-service/     # User management service
│   └── web-service/      # Next.js frontend
├── migrations/           # Database migrations
└── docker-compose.yml    # Docker compose configuration
```

## API Documentation
The API Gateway service exposes RESTful endpoints and WebSocket connections for:
- User authentication
- Chat management
- Message streaming

You can find the API Documentation [here](https://documenter.getpostman.com/view/8647915/2sAYJ6CzwA).

## Contributing
1. Fork the repository
2. Create a new branch (`git checkout -b feature/awesome-feature`)
3. Commit the changes (`git commit -am 'Add awesome feature'`)
4. Push to the branch (`git push origin feature/awesome-feature`)
5. Create a new Pull Request