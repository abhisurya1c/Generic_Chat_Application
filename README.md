# Generic Chat Application

A modern, ChatGPT-style web application for general-purpose AI interaction using a local LLM via Ollama. 
Built with **Golang** (Backend), **React** (Frontend), and **PostgreSQL** (Database).

*(Formerly SQL Assistant Web Application)*

## Features

- **Generic AI Assistant**: Helpful, harmless, and honest assistant for coding, writing, and analysis.
- **User Authentication**: Secure Login and Registration system using JWT.
- **Persistent Chat History**: All chats and messages are saved to a local PostgreSQL database.
- **Real-time Streaming**: Responses are streamed token-by-token using Server-Sent Events (SSE).
- **Chat Management**: Create new chats, view history in sidebar, and delete old conversations.
- **Dark Mode UI**: Clean, responsive interface inspired by modern chat apps.

## Prerequisites

1.  **Go** (1.21+)
2.  **Node.js** (18+)
3.  **Docker & Docker Compose**: Required for the PostgreSQL database.
4.  **Ollama**: Installed and running locally.
    - Pull a model: `ollama pull llama3` (default) or any other supported model.

## Running the Application

### 1. Start the Database
Start the PostgreSQL container using Docker Compose:

```bash
docker-compose up -d
```
This runs Postgres on port 5432 with the credentials defined in `docker-compose.yml`.

### 2. Start the Backend

```bash
cd backend
go mod tidy
go run main.go
```
The server will start at `http://localhost:8080`.
It will automatically connect to the database and run any necessary migrations (create tables).

### 3. Start the Frontend

```bash
cd frontend
npm install
npm run dev
```
The application will be available at `http://localhost:5173`.

## Usage

1.  Register a new account on the login screen.
2.  Log in with your credentials.
3.  Type a message to the AI (e.g., "Write a poem about coding").
4.  Your chat history is saved automatically. You can access previous chats from the sidebar.
5.  To delete a chat, hover over it in the sidebar and click the trash icon.

## Project Structure

- `backend/`: Golang server
    - `handlers/`: API route handlers (Auth, Chat, History).
    - `middleware/`: JWT Authentication and CORS.
    - `db/`: Database connection and schema migrations.
    - `services/`: Ollama API client.
- `frontend/`: React + Vite application
    - `src/components/`: Reusable UI components (Login, Sidebar, ChatWindow).
    - `src/api.js`: API client methods.
- `docker-compose.yml`: PostgreSQL service configuration.

## Environment Variables

- **Backend**:
    - Port: 8080
    - DB Connection: `postgres://user:password@localhost:5432/sqlchat` (Hardcoded in `main.go` for dev)
    - JWT Secret: `my_secret_key` (Hardcoded in `middleware/auth.go` for dev)

- **Frontend**:
    - API URL: `http://localhost:8080/api`