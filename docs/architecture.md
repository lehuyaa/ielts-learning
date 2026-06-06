# Architecture - React.js + Go Gin + MySQL + GORM

## 1. Chosen Stack

```txt
Frontend: React.js + TypeScript + Vite
Backend: Go + Gin
Database: MySQL 8+
ORM: GORM
Auth: JWT
Styling: Tailwind CSS + shadcn/ui
Frontend HTTP client: axios
```

---

## 2. Monorepo Structure

Recommended:

```txt
ielts-vocab-app/
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   └── .env.example
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── modules/
│   │   │   ├── auth/
│   │   │   ├── dashboard/
│   │   │   ├── roadmap/
│   │   │   ├── lesson/
│   │   │   ├── vocabulary/
│   │   │   ├── flashcard/
│   │   │   ├── quiz/
│   │   │   ├── profile/
│   │   │   └── placement/
│   │   └── shared/
│   │       ├── response/
│   │       ├── errors/
│   │       └── validator/
│   ├── seeds/
│   ├── go.mod
│   └── .env.example
│
├── docs/
├── AGENTS.md
├── docker-compose.yml
└── README.md
```

---

## 3. Frontend Architecture

### Frontend Folder Structure

```txt
frontend/src/
├── api/
│   ├── api.ts
│   ├── auth.ts
│   ├── dashboard.ts
│   ├── roadmap.ts
│   ├── lesson.ts
│   ├── vocabulary.ts
│   ├── flashcard.ts
│   ├── quiz.ts
│   └── profile.ts
├── contexts/
│   └── auth/
│       ├── AuthContext.tsx
│       ├── AuthProvider.tsx
│       └── useAuth.ts
├── types/
│   ├── api.ts
│   ├── auth.ts
│   ├── dashboard.ts
│   ├── roadmap.ts
│   ├── lesson.ts
│   ├── vocabulary.ts
│   ├── flashcard.ts
│   ├── quiz.ts
│   └── profile.ts
├── components/
│   ├── ui/
│   ├── layout/
│   └── shared/
│
├── features/
│   ├── auth/
│   ├── dashboard/
│   ├── roadmap/
│   ├── lesson/
│   ├── vocabulary/
│   ├── flashcard/
│   ├── quiz/
│   └── profile/
│
├── hooks/
├── pages/
│   ├── LandingPage.tsx
│   ├── LoginPage.tsx
│   ├── RegisterPage.tsx
│   ├── DashboardPage.tsx
│   ├── RoadmapPage.tsx
│   ├── LessonDetailPage.tsx
│   ├── FlashcardPage.tsx
│   ├── ReviewDuePage.tsx
│   ├── QuizPage.tsx
│   ├── VocabularyListPage.tsx
│   ├── VocabularyDetailPage.tsx
│   └── ProfilePage.tsx
├── lib/
│   ├── utils.ts
│   └── constants.ts
├── main.tsx
└── index.css
```

Rules:

- Use axios for all frontend API calls.
- Do not use `fetch`.
- The shared axios instance must live in `frontend/src/api/api.ts`.
- Feature API modules must live in `frontend/src/api`.
- All React contexts must live in `frontend/src/contexts`.
- Each context must have its own folder, for example `frontend/src/contexts/auth`.
- All shared TypeScript types must live in `frontend/src/types`.
- Do not scatter shared types across feature folders.
- `frontend/src/lib` is for non-API utilities only.

---

## 4. Frontend Routing

Use:

```txt
react-router-dom
```

Routes:

```txt
/
 /login
 /register
 /placement-test
 /dashboard
 /roadmap
 /lessons/:lessonId
 /lessons/:lessonId/flashcards
 /lessons/:lessonId/quiz
 /reviews
 /vocabulary
 /vocabulary/:vocabularyId
 /profile
```

Protected routes:

```txt
/dashboard
/roadmap
/lessons/*
/reviews
/vocabulary/*
/profile
```

Create:

```txt
frontend/src/features/auth/ProtectedRoute.tsx
```

---

## 5. Frontend API Client

Create:

```txt
frontend/src/api/api.ts
```

Responsibilities:

- Set base URL
- Attach JWT token
- Parse standard response shape
- Handle 401 globally
- Export configured axios instance and typed response helpers
- Do not use `fetch`

Example:

```ts
import axios from "axios"

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken")

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})
```

Feature API modules should use the shared instance:

```ts
import { api } from "./api"

export async function login(input: LoginRequest): Promise<AuthResponse> {
  const response = await api.post<{ data: AuthResponse }>("/auth/login", input)

  return response.data.data
}
```

---

## 6. Frontend State Management

For MVP:

- Use React state
- Use React Context for auth
- Use URL search params for filters
- Optional: TanStack Query for server state

Recommendation:

Use TanStack Query if you want cleaner API state handling.

If keeping simple:

- useEffect + useState is acceptable for MVP

Avoid:

- Redux
- MobX
- Complex global stores

Context placement:

- Put contexts in `frontend/src/contexts`.
- Put the auth context in `frontend/src/contexts/auth`.
- Feature folders may import context hooks, but should not define shared contexts.
- Shared request/response/domain types belong in `frontend/src/types`.

---

## 7. UI Source Of Truth

When a page screenshot exists:

```txt
docs/screenshots/<page>.png
```

The screenshot is the authoritative design reference.

The implementation should reproduce as closely as possible:

- Layout
- Section order
- Spacing
- Typography hierarchy
- Component structure

If screenshots and design preferences conflict, follow the screenshot.

---

## 8. Validation Architecture

Validation contracts:

- `docs/validation-contracts.md` is the source of truth for user-facing form validation.
- Every user-facing form must have a validation contract before implementation.
- Frontend and backend validation must both match the documented contract.
- Login/Register forms must use the auth validation contracts.

### Frontend Validation

Use Zod schemas.

Rules:

- Put shared form schemas in a predictable location.
- Recommended location:

```txt
frontend/src/features/<feature>/validation/
```

Examples:

```txt
frontend/src/features/auth/validation/authSchemas.ts
frontend/src/features/profile/validation/profileSchemas.ts
```

- Form schemas must be derived from `docs/validation-contracts.md`.
- Frontend validation is for user experience and fast feedback.
- Frontend validation must not be the only validation layer.

### Backend Validation

Validate request DTOs.

Rules:

- Backend validation must enforce the same rules as `docs/validation-contracts.md`.
- Normalize email, username, and name before use.
- Email normalization: trim and lowercase.
- Username normalization: trim and lowercase.
- Name normalization: trim.
- Return field-level validation errors using the standard `VALIDATION_ERROR` response shape from `docs/api.md`.
- Never rely only on frontend validation.
- Backend validation is the final authority.

---

## 8. Backend Architecture

### Backend Folder Structure

```txt
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── mysql.go
│   │   └── migrate.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/
│   │   └── models.go
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── routes.go
│   │   ├── dashboard/
│   │   ├── roadmap/
│   │   ├── lesson/
│   │   ├── vocabulary/
│   │   ├── flashcard/
│   │   ├── quiz/
│   │   └── profile/
│   └── shared/
│       ├── response/
│       ├── app_error/
│       ├── jwt/
│       └── password/
├── seeds/
│   └── seed.go
├── go.mod
└── .env.example
```

---

## 9. Backend Module Pattern

Each module should have:

```txt
handler.go     HTTP layer
service.go     business logic
repository.go  database queries
dto.go         request/response structs
routes.go      route registration
```

Example:

```txt
internal/modules/flashcard/
├── handler.go
├── service.go
├── repository.go
├── dto.go
└── routes.go
```

---

## 10. Backend Layer Responsibilities

### Handler

Does:

- Read request
- Validate input
- Get user ID from context
- Call service
- Return response

Does not:

- Put business logic
- Write complex DB queries

### Service

Does:

- Business logic
- Transactions
- Progress update
- Quiz scoring
- Spaced repetition rules

### Repository

Does:

- GORM queries
- Create/update/find records

---

## 11. Main Backend Dependencies

Recommended Go packages:

```txt
github.com/gin-gonic/gin
gorm.io/gorm
gorm.io/driver/mysql
github.com/golang-jwt/jwt/v5
golang.org/x/crypto/bcrypt
github.com/joho/godotenv
github.com/go-playground/validator/v10
```

Optional:

```txt
github.com/rs/zerolog
github.com/swaggo/gin-swagger
```

---

## 12. Environment Variables

### Backend `.env`

```txt
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=ielts_vocab

JWT_SECRET=change_me
JWT_EXPIRES_IN_HOURS=168

FRONTEND_URL=http://localhost:5173
```

### Frontend `.env`

```txt
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## 13. Docker Compose

Recommended root `docker-compose.yml`:

```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: ielts_vocab_mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: ielts_vocab
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
```

---

## 14. Backend Route Registration

Example:

```go
func RegisterRoutes(r *gin.Engine, deps Dependencies) {
	api := r.Group("/api/v1")

	auth.RegisterRoutes(api, deps.DB, deps.Config)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(deps.Config.JWTSecret))

	dashboard.RegisterRoutes(protected, deps.DB)
	roadmap.RegisterRoutes(protected, deps.DB)
	lesson.RegisterRoutes(protected, deps.DB)
	vocabulary.RegisterRoutes(protected, deps.DB)
	flashcard.RegisterRoutes(protected, deps.DB)
	quiz.RegisterRoutes(protected, deps.DB)
	profile.RegisterRoutes(protected, deps.DB)
}
```

---

## 15. Auth Architecture

Use JWT.

Login/register returns:

```json
{
  "accessToken": "..."
}
```

JWT payload:

```json
{
  "userId": 1,
  "email": "alex@example.com",
  "role": "USER",
  "exp": 123456789
}
```

Backend middleware:

- Reads Authorization header
- Validates token
- Sets user ID into Gin context

Example context keys:

```txt
userID
userEmail
userRole
```

---

## 16. CORS

Backend should allow frontend origin:

```txt
http://localhost:5173
```

Allow headers:

```txt
Authorization
Content-Type
```

Allow methods:

```txt
GET
POST
PUT
PATCH
DELETE
OPTIONS
```

---

## 17. Service Logic Placement

### Spaced repetition

```txt
backend/internal/modules/flashcard/service.go
```

### Quiz scoring

```txt
backend/internal/modules/quiz/service.go
```

### Lesson completion

```txt
backend/internal/modules/lesson/service.go
```

### Dashboard summary

```txt
backend/internal/modules/dashboard/service.go
```

---

## 18. Development Workflow

Run MySQL:

```bash
docker compose up -d
```

Run backend:

```bash
cd backend
go run ./cmd/api
```

Run frontend:

```bash
cd frontend
pnpm dev
```

---

## 19. Testing Strategy

For MVP:

Backend unit tests:

- Auth service
- Flashcard spaced repetition
- Quiz scoring
- Lesson completion

Frontend manual tests:

- Register/login
- Protected route redirect
- Dashboard load
- Roadmap load
- Flashcard session
- Quiz submit

Later:

- Go integration tests
- Playwright E2E
- API contract tests
