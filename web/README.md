# Web Package

The `web` package provides a simple RESTful API for interacting with the `godb` database over HTTP.

## Import

To use the `web` package, import it as follows:

```go
import "godb/web"
```

## Usage

The main component of the `web` package is the `Server` struct. You can create a new server instance using the `NewServer` function, which takes a port number as input (e.g., `:8080`). The `Start` method starts the server.

**Example:**

To start a new web server on port 8080:

```go
package main

import (
	"godb/web"
	"log"
)

func main() {
	server := web.NewServer(":8080")

	// Initialize database schema
	if err := server.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```

## API Endpoints

The web server exposes the following RESTful endpoints:

### Users

-   `POST /users`: Creates a new user.
    -   **Request Body:**
        ```json
        {
            "id": 1,
            "name": "moses",
            "email": "moses@example.com"
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "User created successfully",
            "count": 1
        }
        ```
-   `GET /users`: Retrieves a list of all users.
    -   **Response:**
        ```json
        [
            {
                "id": 1,
                "name": "moses",
                "email": "moses@example.com"
            }
        ]
        ```

### Posts

-   `POST /posts`: Creates a new post.
    -   **Request Body:**
        ```json
        {
            "id": 101,
            "user_id": 1,
            "title": "My First Post",
            "body": "This is the content of my first post."
        }
        ```
    -   **Response:**
        ```json
        {
            "message": "Post created successfully",
            "count": 1
        }
        ```
-   `GET /posts`: Retrieves a list of all posts, joined with user information.
    -   **Response:**
        ```json
        [
            {
                "post_id": 101,
                "post_title": "My First Post",
                "post_body": "This is the content of my first post.",
                "user_id": 1,
                "user_name": "moses",
                "user_email": "moses@example.com"
            }
        ]
        ```

## Components

### Server

The `Server` struct is responsible for setting up the database, registering the API endpoints, and starting the web server.

### Handler

The `Handler` struct contains the HTTP handlers for the API endpoints. These handlers are responsible for parsing requests, calling the appropriate `engine` methods, and sending back responses.

### Data Transfer Objects (DTOs)

The `dto.go` file defines the Data Transfer Objects (DTOs) for the API. These are the structs that are used to serialize and deserialize JSON requests and responses.