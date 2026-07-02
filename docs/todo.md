# shaw-ty Implementation Checklist

Structured backlog for the URL Shortener API. Tasks are grouped by urgency:

- **Immediate** — fix bugs and close gaps in what is already built so the API works reliably as structured today
- **Intermediate** — spec-aligned features and security polish needed for a complete API
- **Long-term** — deployment, tests, frontend, telemetry, and other resume-project enhancements not strictly required to finish the current backend

> Replaces the former root `todo.md`, which was stale (many items were already implemented).

---

## Immediate

*Fix bugs and close gaps in what is already built so the API works reliably as designed today.*

### Controller & routing bugs (highest priority)

- Fix missing `return` after error responses in `[api/controllers/link-controller.go](../api/controllers/link-controller.go)` (`CreateLink`, `GetLink`, `GetFullLink`, `UpdateLink`) and `[api/controllers/user-controller.go](../api/controllers/user-controller.go)` (`ReactivateUser`)
- Fix `UpdateLink` route mismatch: route is `PUT /short_links` (no `:id`) but handler reads `c.Param("id")` — change to `PUT /short_links/:id` in `[api/routes/routes.go](../api/routes/routes.go)` and add nil checks for optional pointer fields in the DTO
- Fix `GetFullLink` redirect toggle: reads `c.Param("redirect")` but route has no such param — use query param (e.g. `?redirect=false`) or remove dead logic
- Fix `SetActiveById` column name bug in `[internal/storage/db/shortlink-repository.go](../internal/storage/db/shortlink-repository.go)`: `SetColumn("isActive", ...)` should use `is_active` (Bun column name)
- Fix broken test import in `[test/internal/utils_test.go](../test/internal/utils_test.go)` (`URLShortener/utils` → `URLShortener/internal/utils`)

### Auth & session correctness

- Enforce JWT token blacklist in `AuthMiddleware` (`[internal/middleware/jwt.go](../internal/middleware/jwt.go)`) — logout currently writes to blacklist but middleware never checks it
- Reject login for deactivated users — check `IsActive` in `[api/controllers/auth-controller.go](../api/controllers/auth-controller.go)` / `[internal/services/user-service.go](../internal/services/user-service.go)`
- Re-hash password on user update — `[UpdateUser](../api/controllers/user-controller.go)` passes raw password to service without Argon2 hashing
- Handle salt creation error in register — currently `salt, _ := auth.CreateUserSalt(...)` in auth controller

### Core URL shortener behavior

- Add **public** redirect route `GET /:code` (outside `/api/v1`, no JWT) per `[api/api-spec.md](../api/api-spec.md)` — currently redirect only exists on protected `GET /api/v1/short_link/:code`
- Enforce `IsActive` and `ExpirationDate` when resolving links for redirect/get
- Wire `expiresAfter` from `[api/dto/link-dto.go](../api/dto/link-dto.go)` into `CreateLink` (currently hardcoded `0`)
- Clarify custom code vs generated code: service stores `CustomCode` but returns Base63 snowflake code — decide and implement consistent behavior in `[internal/services/link-service.go](../internal/services/link-service.go)`

### Security & data exposure (minimum viable)

- Stop returning full `models.User` (includes `Password`, `Salt`) — use existing `SafeUser` model in `[internal/storage/models/user.go](../internal/storage/models/user.go)` or response DTOs in `[api/dto/user-dto.go](../api/dto/user-dto.go)`
- Replace service-layer `panic(err)` with returned errors in `[internal/services/user-service.go](../internal/services/user-service.go)` and `[internal/services/link-service.go](../internal/services/link-service.go)`
- Sanitize login error responses — stop returning raw DB errors to clients (`[api/controllers/auth-controller.go](../api/controllers/auth-controller.go)`)

### Config & local dev portability

- Use `cfg.ServerHost` / `cfg.ServerPort` in `[cmd/main.go](../cmd/main.go)` instead of hardcoded `:8080`
- Replace machine-specific absolute DB path in `[internal/config/config.yaml](../internal/config/config.yaml)` with a relative/portable default (e.g. `./sql/url_shortener.db`)
- Update `[scripts/setup-db-basic.sh](../scripts/setup-db-basic.sh)` paths to match current repo location
- Move machine ID validation from `[cmd/main.go](../cmd/main.go)` into `[internal/config/config.go](../internal/config/config.go)` load (per existing TODO)
- Call `c.Abort()` in rate limiter when limit exceeded (`[internal/middleware/rate-limit.go](../internal/middleware/rate-limit.go)`)
- Read JWT expiry from config (`jwt_expiry_duration` in yaml) instead of hardcoded 24h in `[internal/middleware/jwt.go](../internal/middleware/jwt.go)`

---

## Intermediate

*Features implied by the existing design, `[api/api-spec.md](../api/api-spec.md)`, and inline TODOs — needed for a complete, spec-aligned API but not blocking basic local use.*

### Authorization & access control

- Enforce ownership: authenticated users can only read/update/delete their own links and profile (TODOs in all three controllers)
- Restrict `GET /user` (list all users) to admin role or remove from non-admin access
- Implement RBAC/RBAP rules in services (5+ TODOs in `[internal/services/user-service.go](../internal/services/user-service.go)`)
- Add permission checks on `GetAllLinksByUser` — currently any user can query any `user_id` query param

### Error handling & API responses

- Introduce dev vs production error responses — use `is_dev_mode` in `[internal/config/config.yaml](../internal/config/config.yaml)` to return descriptive errors (stack traces, underlying causes) in dev and generic/safe messages in prod
- Propagate typed/sentinel errors from services and map them to HTTP status codes in controllers (replace ad-hoc string errors and inconsistent shapes)
- Standardize error response shape (`{error}`, `{status, error}`, `{success, message}` are inconsistent today)
- Generify error messages in production — avoid leaking internal details (OWASP TODOs in auth/link controllers); keep detailed messages server-side only (structured logs)

### Input validation & API polish

- Add input sanitization and validation across controllers (URL format, email, password strength, username rules)
- Fix `didAnyChange` logic bug in `[internal/services/user-service.go](../internal/services/user-service.go)` (`UpdateCompleteUserById` sets wrong field for `lastPasswordUpdate`)
- Align routes with spec where intentional divergence exists (e.g. `POST /register` vs spec's `POST /user`, `GET /short_links?user_id=` vs spec's `GET /user/:id/links`)

### Link & user feature completeness

- Add pagination and date filtering for user links (spec: `/user/{id}/links`)
- Implement click/analytics counter on redirect (DTO has `Clicks` field; no logic yet)
- Idempotency enforcement for register, link create, and link update (spec requirement)
- Password expiry signal on user reactivation (TODO in user controller/service)
- Populate dev seed data in `[sql/seeding/seed.go](../sql/seeding/seed.go)` (scaffold exists, `usersToSeed` is empty)
- Remove or externalize hardcoded admin password in `[sql/migrations/20250609_initial_population.go](../sql/migrations/20250609_initial_population.go)`

### i18n & DTO layer

- Migrate hardcoded controller strings to `[internal/i18n/locales/en.yaml](../internal/i18n/locales/en.yaml)`
- Move model→DTO conversion from controllers into services (TODO in user/link controllers)
- Expand i18n beyond the single existing error key

### Testing (core workflows)

- Add unit tests for auth (register, login, logout, token validation)
- Add unit tests for link service (create, get, redirect logic, expiration)
- Add controller tests with `httptest` for happy-path and error cases
- Fill in stub test cases in `[test/config/config_test.go](../test/config/config_test.go)`

---

## Long-term

*Not required to finish the backend as structured today, but would meaningfully evolve this into a polished resume/portfolio project.*

### Containerization & deployment

- Add multi-stage `Dockerfile` (non-root user, static binary)
- Add `docker-compose.yml` (app + Postgres for prod-like local dev)
- Add `.env.example` documenting required env vars (`APP_PEPPER`, `BASE_URL`, `JWT_SECRET`, DB vars)
- Add `/healthz` and `/readyz` endpoints for load balancers
- Implement graceful shutdown (signal handling in `[cmd/main.go](../cmd/main.go)`)
- Deploy to a cloud target (AWS ECS/Fargate, Fly.io, Railway, etc.) — backend + optional frontend

### Kubernetes & production hardening

- K8s manifests or Helm chart (Deployment, Service, Ingress, ConfigMap, Secret)
- Switch production DB from SQLite to Postgres — add driver to `[go.mod](../go.mod)`, finish wiring in `[internal/storage/db/go-dbc.go](../internal/storage/db/go-dbc.go)` (Postgres/MySQL stubs exist but drivers missing)
- Externalize in-memory token blacklist to Redis (required for horizontal scaling)
- Externalize rate limiter state (Redis) for multi-replica deployments
- Add CORS middleware for cross-origin frontend
- Secrets management (K8s Secrets, Vault, or cloud SM) — remove `.env`-only pattern for prod
- TLS termination (Ingress + cert-manager or cloud LB)
- Prod migration strategy (job/init container vs startup migrate)

### CI/CD & quality gates

- GitHub Actions: `go test`, `go vet`, `staticcheck`, build
- CI: build and push Docker image on merge/tag
- Test coverage reporting (README still says "TODO: coverage")
- Lint/format enforcement (`golangci-lint`, `gofumpt`)
- Postman collection or OpenAPI/Swagger spec generated from routes

### Observability & telemetry

- Structured JSON logging (replace ad-hoc `log.Printf`)
- Prometheus metrics (request latency, redirect count, auth failures)
- OpenTelemetry tracing for request → service → DB
- Alerting on SLO violations (spec: login <500ms, redirect <100ms server time)

### Frontend & client

- Build SPA (React/Vue/Svelte): login, shorten form, link list
- JWT storage and refresh handling on client
- Public short-link landing/redirect page at root domain
- Optional: CLI client or Go SDK for API consumers

### Microservices & architecture

- Extract auth into a standalone gRPC microservice (register, login, token issue/validate, blacklist) — main API becomes a gRPC client for auth instead of in-process JWT middleware
- Define protobuf contracts for auth RPCs and document service boundaries (what stays in the monolith vs moves out)
- Add inter-service auth (mTLS or shared secret) and health checks for the auth service in compose/K8s

### Refactors & tech debt

- Migrate from deprecated `github.com/dgrijalva/jwt-go` to `github.com/golang-jwt/jwt`
- Replace runtime `panic` with returned errors in services, controllers, and libraries (`[internal/utils/utils.go](../internal/utils/utils.go)`, `[internal/core/encoder/base63.go](../internal/core/encoder/base63.go)`) — reserve `log.Fatal` for unrecoverable startup failures only (missing config, DB connect, migration failure)
- Audit all `panic`/`log.Fatal` call sites and document the convention: `log.Fatal` at process init, returned `error` on request path, never panic in handlers
- Consolidate duplicate update paths in user service (`UpdateUserById` vs `UpdateCompleteUserById`)
- Add repository interface mocks for service unit tests
- Consider admin role model in DB schema (currently no roles table)

### Future features (aspirational)

- Google/OAuth login
- Analytics dashboard (click trends, top links)
- Link expiration UI and email notifications
- User profile management page
- API rate limits per user (not just global IP-based)
- Custom domain support for short links

---

## Notes

- Tackle **Immediate** tasks in listed order: bugs first, then auth, then redirect, then config
- Each checkbox is scoped to a single PR-sized unit where possible
- See `[api/api-spec.md](../api/api-spec.md)` for functional requirements and `[DATABASE_SETUP.md](../DATABASE_SETUP.md)` for local DB setup

