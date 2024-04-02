# golang-api

## Instructions

### Running this server

1. Create a PostgreSQL database for this API server. Then, start the database and run `schema.sql` to create the necessary tables for the server.
2. Create a `.env` file in the root directory of this project, with the same fields as those in `.env.example`. Include the information used to set up your PostgreSQL database for fields starting with `DB_`.
3. `cd` into this project's root directory and run `go run ./cmd/main.go` in the terminal to start this project.

### Testing

1. Test the APIs with this Postman [collection](https://www.postman.com/security-physicist-88556430/workspace/shared/collection/28021553-357d524d-b51c-47d8-8054-fea64ddacd09?action=share&creator=28021553). You may also create your own Postman collection for testing.
