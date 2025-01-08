# Instant Messaging App

An instant messaging application with a modern architecture, featuring user authentication, message exchanges, and group chat functionalities. The app leverages RabbitMQ for message communication, PostgreSQL for data persistence, and a microservices architecture for scalability.

## Features

- Real-time messaging with RabbitMQ.
- User authentication and JWT-based session management.
- Private messaging.
- Scalable microservices architecture.
- Frontend built with Vite.js, React, and Tailwind CSS.
- Backend microservices built with Go (Fiber framework).
- Database: PostgresSQL.
- Message broker: RabbitMQ.

## Tech Stack

### Frontend

- Framework: Vite.js + React
- Styling: Tailwind CSS
- Language: TypeScript

### Backend

- Language: Go (Golang)
- Framework: Fiber
- Authentication: JWT

### Infrastructure

- Database: PostgreSQL
- Message Broker: RabbitMQ
- Web Server: Caddy

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js (for local frontend development)
- Go (for local backend development)

### Installation

1. Clone the repository:

```
git clone https://github.com/tchadelicard/instant-messaging-app
cd instant-messaging-app
```

2. Start the application with Docker Compose:

```
docker compose up -d
```

4. Access the app:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- RabbitMQ Management: http://localhost:15672 (Default credentials: guest / guest)

## Project structure

```
.
├── api                   # Backend API services
├── cmd                   # Command-line entry points
├── config                # Configuration code
├── frontend              # Frontend (React + Vite.js)
│   ├── src
│   │   ├── components    # Reusable React components
│   │   ├── pages         # Page-specific components
│   │   ├── services      # API integrations
│   │   └── types         # TypeScript interfaces
├── models                # Database models
├── user                  # User-related services
├── message               # Message-related services
├── utils                 # Utility functions
├── docker-compose.yml    # Docker Compose file
├── Dockerfile            # Backend Dockerfile
├── init.sql              # Database initialization script
└── README.md             # README
```

## Usage without Docker

1. Frontend:

Start locally:

```
cd frontend
npm install
npm run dev
```

Build production version:

```
npm run build
```

2. Backend:

Run the backend locally:

```
go run main.go api
go run main.go user
go run main.go message
```

3. Database:

```
psql -h localhost -U postgres -d instant_messaging_app
```

4. RabbitMQ:

Open the RabbitMQ management UI at http://localhost:15672.

## Environment variables

| Variable            | Description              | Default                 |
| ------------------- | ------------------------ | ----------------------- |
| `DB_HOST`           | PostgreSQL host          | `postgres`              |
| `DB_USER`           | PostgreSQL username      | `postgres`              |
| `DB_PASSWORD`       | PostgreSQL password      | `postgres`              |
| `DB_NAME`           | PostgreSQL database name | `instant_messaging_app` |
| `DB_PORT`           | PostgreSQL port          | `5432`                  |
| `JWT_SECRET`        | Secret key for JWT       | `your-secret-key`       |
| `APP_PORT`          | Application port         | `8080`                  |
| `RABBITMQ_HOST`     | RabbitMQ host            | `rabbitmq`              |
| `RABBITMQ_PORT`     | RabbitMQ port            | `5672`                  |
| `RABBITMQ_USER`     | RabbitMQ username        | `guest`                 |
| `RABBITMQ_PASSWORD` | RabbitMQ password        | `guest`                 |
