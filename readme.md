# Simplified Vault-like Service

This is a simplified implementation of a Vault-like service in Go. The service allows you to create roots with secrets, create users, AppRoles, generate tokens, and fetch secrets using tokens.

## Features

- Create roots with secrets
- Create users
- Create AppRoles
- Generate tokens
- Fetch secrets using tokens

## Getting Started

### Prerequisites

- Go 1.16 or later
- cURL

### Running the Application

1. Clone the repository:

    ```sh
    git clone https://github.com/your-repo/vault-like-service.git
    cd vault-like-service
    ```

2. Build and run the application:

    ```sh
    go build -o vault-like-service
    ./vault-like-service
    ```

3. The service will be running on `http://localhost:8080`.

### API Endpoints

- `POST /create-root?name={name}&value={value}&ttl={ttl}`: Create a root with a secret.
- `POST /create-user?username={username}&password={password}`: Create a user.
- `POST /create-approle?roleID={roleID}&secretID={secretID}`: Create an AppRole.
- `GET /get-token?roleID={roleID}&secretID={secretID}`: Generate a token using AppRole.
- `GET /get-secret?token={token}&rootName={rootName}`: Fetch a secret using a token.

### Testing with cURL

1. **Create a Root with a Secret**

    ```sh
    curl -X POST "http://localhost:8080/create-root?name=mysecret&value=supersecret&ttl=1h"
    ```

2. **Create a User**

    ```sh
    curl -X POST "http://localhost:8080/create-user?username=user1&password=password123"
    ```

3. **Create an AppRole**

    ```sh
    curl -X POST "http://localhost:8080/create-approle?roleID=approle1&secretID=secret123"
    ```

4. **Get a Token using the AppRole**

    ```sh
    TOKEN=$(curl -s -X GET "http://localhost:8080/get-token?roleID=approle1&secretID=secret123")
    echo "Generated Token: $TOKEN"
    ```

5. **Fetch a Secret using the Token**

    ```sh
    curl -X GET "http://localhost:8080/get-secret?token=$TOKEN&rootName=mysecret"
    ```

