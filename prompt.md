# Generic_Chat_Application

Prompt :
You are a senior full-stack engineer and UI/UX designer.

Your task is to generate a complete ChatGPT-style SQL Assistant web application using Golang as the backend.

========================
TECH STACK
========================
Frontend:
- HTML5, CSS3, JavaScript (ES6)
- No heavy framework (Vanilla JS preferred)
- Responsive UI inspired by ChatGPT / Claude / Perplexity

Backend:
- Golang (net/http)
- REST APIs
- Server-Sent Events (SSE) for streaming
- CORS enabled

========================
APPLICATION GOALS
========================
1. ChatGPT-Like Interface:
   - Chat panel centered
   - User messages on right, assistant on left
   - Typing animation
   - Auto scroll
   - Input text area with Send button
   - Enter key submits message
   - Code blocks styled for SQL output with copy button

2. Chat History:
   - Sidebar showing:
     - "New Chat"
     - Previous chat sessions
   - Store chat history in browser localStorage
   - Clicking previous chat restores messages
   - Persist data on page reload

3. SQL Query Generator:
   - Assistant returns ONLY SQL by default
   - Explanation only when user explicitly asks
   - Support PostgreSQL, MySQL, SQL Server
   - Add a system prompt before every request:
     "You are an expert SQL assistant. Return ONLY valid SQL queries."

4. Backend API Integration (IMPORTANT):
   - Integrate the following model API using Go HTTP client:

     curl --location 'http://localhost:11434/api/generate' \
     --header 'Content-Type: application/json' \
     --data '{
        "model": "llama3 or sqlcoder",
        "prompt": "<USER_PROMPT>",
        "stream": false
     }'

   - Backend exposes:
     POST /api/chat        → non-stream response
     GET  /api/chat/stream → streaming response (SSE)

   - Frontend NEVER calls the model API directly
   - Backend injects system prompt + user prompt
   - Backend formats response before sending

5. Streaming Support:
   - When stream=true:
     - Use Server-Sent Events (SSE)
     - Stream tokens like ChatGPT
     - Frontend renders text progressively
   - When stream=false:
     - Return full response at once

6. UI/UX:
   - Dark mode ChatGPT-style design
   - Smooth animations
   - Loading indicator while waiting
   - Sidebar collapses on mobile
   - Modern typography

7. Project Structure:
   /frontend
     index.html
     styles.css
     app.js

   /backend
     main.go
     handlers/chat.go
     services/ollama.go
     middleware/cors.go

8. Code Quality:
   - Clean, idiomatic Go
   - Context usage
   - Proper error handling
   - No global state misuse
   - Commented code

========================
SYSTEM PROMPT (MANDATORY)
========================
Always prepend this system message to user input:

You are an expert SQL assistant. 
Your sole purpose is to generate valid, optimized, and secure SQL queries for PostgreSQL, MySQL, and SQL Server.

STRICT OPERATIONAL RULES:
1. OUTPUT FORMAT: Return ONLY valid SQL code. Do NOT include explanations, introductory text, or concluding remarks unless the user explicitly asks for an explanation.
2. SAFETY GUARDRAILS: You are strictly prohibited from generating queries that contain destructive or administrative commands. 
   - FORBIDDEN KEYWORDS: DROP, DELETE, TRUNCATE, ALTER, GRANT, REVOKE, SHUTDOWN.
   - If a user requests a destructive operation, respond with: "-- Error: Destructive SQL commands (DROP, DELETE, TRUNCATE) are not permitted for safety reasons."
3. CODE QUALITY: Ensure all table and column names are handled as described by the user. If schema information is provided, adhere to it strictly.
4. SYSTEM CONTEXT: You are an expert SQL assistant. Return ONLY valid SQL queries. Do NOT explain unless explicitly asked.

========================
OUTPUT EXPECTATION
========================
- Generate full frontend code
- Generate Golang backend code
- Streaming + non-streaming supported
- Provide run instructions
- App should behave like ChatGPT but specialized for SQL generation

Begin generating the project now.