💬 Chat Microservice


Login page
<img width="544" height="378" alt="Screenshot 2026-03-13 083016" src="https://github.com/user-attachments/assets/58136812-0209-481b-ae16-bd1c3ce73578" />

Chat window 
<img width="554" height="629" alt="Screenshot 2026-03-13 083225" src="https://github.com/user-attachments/assets/aaea1462-e33f-44db-8c58-8ce87118683d" />


A simple real-time chat backend built using Go and a microservices approach. The project includes separate services for authentication and chat, along with a basic frontend for testing.

📁 Project Structure
chat-microservice/
│
├── auth-service/        # Authentication service (Go)
├── chat-service/        # Chat service (Go)
├── docker-compose.yml   # Docker configuration
├── index.html           # Frontend for testing chat
⚙️ Prerequisites

Make sure the following are installed:

Go

Docker

Docker Compose

🐳 Running with Docker

Run the following commands from the root directory:

docker compose down
docker compose build --no-cache
docker compose up
🖥️ Running Without Docker
1. Install dependencies
go mod tidy
2. Run Authentication Service
cd auth-service
go run main.go
3. Run Chat Service
cd chat-service
go run main.go
🌐 Running the Frontend

Open index.html in a browser

Open it in two tabs

Use both tabs to test chat functionality

📌 Notes

Ensure required ports are available before running services

Run both services before opening the frontend

Docker must be running for Docker-based execution

🧪 Testing

Open multiple tabs of index.html

Send messages between tabs

Verify real-time communication

🛠️ Tech Used

Go

Docker

HTML (for testing interface)
