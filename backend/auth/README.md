# Auth service

The authentication microservice handles user registration, login, logout, and token refresh operations using the
Keycloack

## HTTP Handlers

- User registration
- User login/logout
- Refreshing token

The service uses the Keycloack for user authentication and authorization. The keycloack realm config in [
`import/trainee-match-realm.json`](import/trainee-match-realm.json)
Keycloack realm has 2 roles (Candidate and Company), so when you register user you choose one of Roles

JWT stores in request cookies with HttpOnly true

## API Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Authenticate user and return tokens
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout user and invalidate tokens

## Configuration

The service can be configured using environment variables:

| Variable                 | Description                              |
|--------------------------|------------------------------------------|
| KC.URL                   | Keycloak server URL                      |
| KC.REALM                 | Keycloak realm                           |
| KC.CLIENT_ID             | OAuth2 client ID                         |
| KC.CLIENT_SECRET         | OAuth2 client secret                     |
| KC.ADMIN_USERNAME        | Keycloak admin username                  |
| KC.ADMIN_PASSWORD        | Keycloak admin password                  |
| KC.ACCESS_TOKEN_EXPIRES  | Access token expiration time in seconds  |
| KC.REFRESH_TOKEN_EXPIRES | Refresh token expiration time in seconds |
| ADDR                     | Service listening address                |

## Running the Service

Please ensure that you have running KeyCloack instance

To start service run Dockerfile

```bash
docker build -t auth_service
docker run -p 8000:8000 auth-service
```