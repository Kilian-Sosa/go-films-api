# 🎬 Go Films API

A film management REST API developed in **Go** using **Gin** and **GORM**, with **MySQL** as the database. The project supports basic film management features such as creating, retrieving, updating, and deleting films. Each film is linked to a registered user (the creator).

This project was built as part of a technical interview to demonstrate backend development skills using **Go**.

---

## 📋 Features

✅ User registration and login (with hashed passwords)  
✅ JWT-based authentication  
✅ Film management (CRUD operations)  
✅ Only the creator can edit or delete a film  
✅ Filtering films by title, genre, and release date  
✅ Full Swagger documentation (OpenAPI 3.0)  
✅ Follows clean architecture (handler, service, repository)  
✅ Docker support (API + MySQL)  
✅ SQL Injection prevention via parameterized queries

---

## 📚 Swagger API Documentation

This project uses **Swagger (swaggo)** to generate and serve API documentation. Once the application is running, you can easily explore and test the endpoints directly from your browser.

### View Swagger UI

After starting the server, visit:

➡️ **http://localhost:8080/swagger/index.html**

This interactive UI allows you to:
- See all available endpoints.
- Inspect request and response formats.
- Execute requests directly from the browser.

### Regenerating Swagger Docs

After modifying any endpoint comments or adding new endpoints, regenerate the docs by running:
```bash
swag init
```
This will update the `docs/` folder with fresh OpenAPI definitions.

---

## 📂 Project Structure

```
.
├── cmd                    # Application entry points
│   ├── server              # Main server (API)
│   └── migrate             # Migration runner
├── internal
│   ├── delivery
│   │   ├── http              # Handlers
│   ├── domain                 # Entities (User, Film)
│   ├── repository              # Database access layer
│   ├── usecase                  # Business logic layer
├── migrations                 # SQL schema & seed data
├── docs                        # Auto-generated Swagger docs
├── Dockerfile                  # Docker build
├── docker-compose.yml          # Docker Compose for API + MySQL
├── README.md                    # This file
├── go.mod                       # Go module
├── go.sum                       # Dependency lockfile
```

---

## 🔧 Code Formatting

This project follows Go's recommended formatting conventions. To ensure consistent formatting across all files, **`goimports`** is used.

### What is `goimports`?
`goimports` is a popular tool that:
- Formats Go code (like `gofmt` does)
- Automatically manages and organizes imports (removing unused imports and adding missing ones)

### Installing `goimports`
If you don't have it installed, you can get it with:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

### Using `goimports`
You can format your files manually with:
```bash
goimports -w ./
```
This will recursively format all files in the project.

### VSCode Auto-Format on Save
If you're using VSCode, configure it to run `goimports` automatically on save:

1. Open settings (Ctrl + ,)
2. Search for `go.formatTool`
3. Set it to:
```json
"go.formatTool": "goimports"
```
4. Ensure format-on-save is enabled:
```json
"editor.formatOnSave": true
```

This ensures your code stays clean and properly formatted with minimal effort.

---

## 🚀 Getting Started

### Architecture Overview

```text
+---------------------+
|   HTTP Handlers     | <--- Gin handles incoming requests
+---------------------+
             |
+---------------------+
|   Usecase Layer     | <--- Business logic (FilmService, UserService)
+---------------------+
             |
+---------------------+
| Repository Layer    | <--- Data access (GORM)
+---------------------+
             |
+---------------------+
|       MySQL         | <--- Persistent storage
+---------------------+
```

---

### Installation using Docker

#### 1. Clone the Repository

```bash
git clone https://github.com/Kilian-Sosa/go-films-api.git
cd go-films-api
```

#### 2. Create .env File

Create a `.env` file at the project root:
```env
DB_HOST=db
DB_USER=root
DB_PASS=root
DB_NAME=database
DB_PORT=3306
APP_PORT=8080
JWT_SECRET=some-secret

MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=database
```

#### 3. Build & Run with Docker

```bash
docker-compose up --build
```

This will spin up:
- `go-films-api` (on port **8080**)
- `go-films-db` (MySQL on port **3306**)

#### 4. Swagger Documentation

Once running, access:
```
http://localhost:8080/swagger/index.html
```

This provides a full, interactive API documentation where you can test requests directly.

---

## 📊 Endpoints

| Method | Endpoint          | Description                    |
|-------|----------------|----------------|
| POST   | `/register`     | Create new user |
| POST   | `/login`        | Login and get token |
| POST   | `/films`        | Create film |
| GET    | `/films`        | List films with filters |
| GET    | `/films/:id`    | Get film details |
| PUT    | `/films/:id`    | Update film (creator only) |
| DELETE | `/films/:id`    | Delete film (creator only) |

---

## 🏆 Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john123","password":"Secret@123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john123","password":"Secret@123"}'
```

### Create Film
```bash
curl -X POST http://localhost:8080/films \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Cool Film",
    "director": "Cool Director",
    "release_date": "2023-04-22",
    "cast": "Actor One, Actor Two",
    "genre": "Drama",
    "synopsis": "A very cool film."
  }'
```

---

## 🛠️ Tech Stack

| Layer       | Technology |
|-------------|-------------|
| Language    | Go 1.24 |
| Framework   | Gin |
| ORM         | GORM |
| Database    | MySQL |
| Auth        | JWT |
| Docs        | Swagger (swaggo) |
| Formatter   | goimports |
| Container   | Docker |
| Tests       | Testify |

---

## 📊 Tests

### Running All Tests
```bash
go test ./...
```

### Test Coverage
- **Unit Tests** for each service (business logic).
- **Repository Tests** using in-memory DB or Docker container.
- **Handler Tests** using Gin's `httptest` package.

### Example Test Command
```bash
go test ./internal/... -v
```

### Mocking
The repository layer is fully mocked in the service-level tests using `testify/mock`.

### Sample Test Output
```
=== RUN   TestRegisterUser_Success
--- PASS: TestRegisterUser_Success (0.02s)
=== RUN   TestLogin_InvalidPassword
--- PASS: TestLogin_InvalidPassword (0.01s)
PASS
ok   	go-films-api/internal/usecase 0.05s
```

---

## 🏅 License

This project is licensed under the **MIT License**.

---

## 💌 Contact

Open an Issue in the repository if you have questions or feedback.

