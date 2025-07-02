# URL Shortener API - TODO List

## 📌 Core API Development
- [ ] Complete `login` endpoint
    - [ ] Implement user authentication
    - [ ] Return JWT token upon successful login
    - [ ] Secure password storage (hashing)

- [ ] Complete `users` endpoints
    - [ ] Implement user registration (`POST /user`)
    - [ ] Implement user retrieval (`GET /user/:id`, `GET /user`)
    - [ ] Implement user update (`PUT /user`)
    - [ ] Implement user deletion (`DELETE /user`)

- [ ] Complete `links` endpoints
    - [ ] Implement link creation (`POST /short_link`)
    - [ ] Implement link retrieval (`GET /short_link/:id`)
    - [ ] Implement retrieving all links (`GET /short_link`)
    - [ ] Implement link update (`PUT /short_link`)
    - [ ] Implement link deletion (`DELETE /short_link`)

- [ ] Implement link redirection
    - [ ] Redirect `GET /:shortCode` to full URL
    - [ ] Handle expired or deactivated links
    - [ ] Implement analytics tracking for redirects (optional)

---

## 🔐 Authentication & Security
- [ ] Implement JWT authentication
    - [ ] Generate JWT on login
    - [ ] Validate JWT in protected routes
    - [ ] Store JWT securely on the frontend
    - [ ] Implement token expiration handling

- [ ] Middleware for authentication & authorization
    - [ ] Secure endpoints with authentication middleware
    - [ ] Ensure users can only access their own links
    - [ ] Implement role-based access control (optional)

---

## 🧪 Testing
- [ ] Write unit tests for:
    - [ ] User authentication
    - [ ] Link creation & retrieval
    - [ ] Token validation
    - [ ] Redirection logic

- [ ] Write integration tests for:
    - [ ] Full user workflows (registration, login, link creation)
    - [ ] Expired or invalid tokens
    - [ ] Unauthorized access attempts

- [ ] Test full API workflow using Postman
    - [ ] User signup & login
    - [ ] Link creation & retrieval
    - [ ] JWT-protected endpoints
    - [ ] Redirection functionality

---

## 🌐 Frontend Development
- [ ] Build a simple frontend UI
    - [ ] Login page
    - [ ] URL shortening form
    - [ ] Display user's shortened URLs

- [ ] Integrate frontend with API
    - [ ] Handle authentication (store JWT)
    - [ ] Make API calls for creating & fetching links
    - [ ] Implement redirecting functionality

---

## 🚀 Deployment & Enhancements
- [ ] Dockerize the application
- [ ] Set up a database for storing users & links
- [ ] Deploy backend & frontend (e.g., AWS, Vercel, or Netlify)
- [ ] Implement analytics (click tracking, user activity logs)
- [ ] Improve error handling & logging

---

### 📌 Notes
- Prioritize **API stability & security** before frontend.
- **Testing is critical** for JWT authentication & link redirection.
- Future features: **link expiration, analytics dashboard, user profiles**.

---

- implement user controller crud
- implement link controller crud
- implement DB
- Test register/login
- test create link
- write base sql/DDL
- obfuscate user info
- google login
- redirect to full url
