# Functional Requirements
- User
  - register a new user
    - idempotent
  - login a user
  - delete the user
  - update a user's information
    - idempotent
  - get the links created by a user
    - limits, pagination, filters?
- Links
  - create a new link
    - idempotent for user
  - delete a link
  - update the link (custom code)
    - idempotent
  - get the full URL given a shortened URL

# Nonfunctional Requirements
- User
  - login should take <500ms
  - all link retrieval should take <500ms
- Links
  - link redirect should take as little time as possible, <100ms my system time + redirect time

# Endpoints
- /user - handles retrieval/edit of person data - requires bearer token
  - GET - takes an id, yields one employee
  - GET - no id yields all employees
  - POST - registers a user, idempotent
  - PUT - updates a user's information, idempotent
  - DELETE - marks the user's account deleted (param to force removal from DB)
  - /{id}/links 
    - GET - returns all the links a user has created (paginated, filtering by dates)
- /login - logs a user in
  - POST - submits USER credentials for login, idempotent within session
  - All other methods rejected
- /short_link - handles retrieval/modification of link data
  - GET - takes an id or a short_code, sends back the full link
  - GET - no parameter, sends back all short links
  - POST - creates a new link
  - PUT - updates a link with a custom code, idempotent
  - DELETE - marks a link inactive
  - TODO: handle redirect here

# Other Information
- Versioning
  - {base_url}/api/v1/{endpoint}
- links
  - https://short.link/{code} is hardcoded to GET