# URL Shortener API

## Features
- Authentication
  - Argon2 auth implementation
  - JWT for sessions
  - timed hash comparison
- Implementation
  - Snowflake algorithm for UUIDs
  - unique, Base63 encoded short links
- Development
  - .env for secrets
  - repository pattern for data layer
  - context'd service calls and error handling logic
- Testing
  - TODO: coverage
- Deployment
  - Docker
  - K8s

## Architecture
Here is a rough breakdown of the architecture of the program, the role each
component plays, and how they interact with each other.

### API Layer
- Router
  - the Router (routes.go) serves the API endpoints through Gin
  - connects the controllers to the endpoints
- Controllers
  - The controllers handle the API logic:
    - Parsing parameters and request bodies
    - serving http responses

### Business Logic Layer
- Services
  - User Service
    - handles user management and CRUD operations
    - transforms data from controllers to models for Data layer operations
    - error handling happens here

### Data Layer
- Repositories
  - The repositories handle the data layer operations, and abstract away the SQL/no-SQL DB operations while handling the data modification requests
  - delivers data modification requests to the database
  - leverages bun ORM around a custom DB controller for abstraction

### Misc
- Middleware
  - JWT for session management after login
- Internal algorithms
  - BASE63 for id > user-facing string conversion
  - Snowflake algorithm for UUIDs
- Auth
  - custom JWT auth implementation around Argon2 encryption
  - salting and peppering
- Config
  - enables application configuration via config.yaml