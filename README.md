# Late Slip Generator

This project is a web application designed to manage student late slips. It provides separate interfaces and functionalities for students and administrators.

## Features

- **User Authentication:** Secure login and registration for both students and admins.
- **Role-Based Access Control:** Differentiates functionalities based on user roles (student, admin).
- **Student Features:**
  - Request late slips.
  - Receive real-time notifications (e.g., approval/rejection of requests) via WebSockets.
- **Admin Features:**
  - Approve or reject late slip requests.
  - View all late slips (pending, approved, rejected).
  - Upload student data.
  - Upload schedule data.
  - Receive real-time notifications via WebSockets.
- **Scheduled Notifications:** A background process (`events.StartScheduleNotifier()`) handles scheduled events or reminders.

## Tech Stack

- **Backend:** Go (Golang)
- **Framework:** Gin (HTTP web framework)
- **Real-time Communication:** WebSockets

## Project Structure (Inferred from `main.go`)

- `controllers/`: Handles incoming HTTP requests and business logic.
- `events/`: Manages real-time event handling (WebSockets, scheduled notifications).
- `initialializers/`: (Likely `initializers/`) Handles application initialization tasks like loading environment variables and connecting to the database.
- `middleware/`: Contains middleware functions for request processing (e.g., authentication, request ID generation, role checking).
- `main.go`: The entry point of the application, sets up routes and starts the server.

## API Endpoints

### General

- `POST /student/register`: Register a new student.
- `POST /student/login`: Log in as a student.
- `POST /admin/register`: Register a new admin.
- `POST /admin/login`: Log in as an admin.

### Student (Requires Authentication & 'student' role)

- `POST /student/requestLateSlip`: Submit a late slip request.
- `GET /student/ws`: Establish a WebSocket connection for real-time updates.

### Admin (Requires Authentication & 'admin' role)

- `PUT /admin/lateslips/approve`: Approve a late slip.
- `GET /admin/lateslips`: Get all late slips.
- `POST /admin/uploadStudentData`: Upload student data.
- `GET /admin/lateslips/pending`: Get all pending late slips.
- `PUT /admin/lateslips/reject`: Reject a late slip.
- `POST /admin/uploadScheduleData`: Upload schedule data.
- `GET /admin/ws`: Establish a WebSocket connection for real-time updates.

## Setup and Running

1.  **Prerequisites:**
    - Go (version specified in `go.mod` if present, otherwise latest stable recommended).
    - A database (details should be in `initialializers/ConnectToDB()` or environment variables).
2.  **Environment Variables:**
    - Create a `.env` file (or configure environment variables as per `initialializers.LoadEnvVariables()`). This will likely include database connection details, JWT secrets, etc.
3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Run the Application:**
    ```bash
    go run main.go
    ```
    The application will start on port `:8000` by default.

---
