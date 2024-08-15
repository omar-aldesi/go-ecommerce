# 🛒 Go Ecommerce API
Welcome to the Go Ecommerce API project! This repository contains a fully-functional RESTful API built using Go, Gin, and Gorm. The project leverages PostgreSQL as its database and includes full Docker integration for seamless development and deployment. Below you'll find all the information you need to clone, set up, and run this project, as well as an overview of its architecture and features.

🎯 Features
Go Gin Gorm: A powerful combination for building RESTful APIs with performance and ease of use.

JWT Authentication: Secure your API endpoints with JSON Web Tokens (JWT) for user authentication.

Custom Error Handling: Robust error handling including validation and type-specific errors.

Modern Project Structure: Clean and scalable code organization with clear separation between concerns.

PostgreSQL Database: Reliable and scalable relational database integration using Gorm.

Full Docker Integration: Effortless development and deployment with Docker, including Docker Compose.

Advanced Currency Handling: Manage transactions and prices effectively with the gocurrency package.

OpenAPI Documentation: Automatically generated API documentation for all endpoints.

RESTful API Protocols: Adheres to REST principles for easy and predictable API interactions.

Auth Middleware: Secure endpoints with a custom JWT-based authentication middleware.

Token Blacklisting: Refresh token blacklisting for enhanced security.

General Systems for Shipping & Payment: Modular and extensible design for handling shipping and payment logic.

Advanced Product Filtering: Powerful and flexible filtering options for product listings.

📁 Project Structure
The project follows a modern Go project structure:

```bash
ecommerce/
├── app/
│   ├── core/           # Core functionality like database and middleware
│   ├── models/         # Data models representing the database schema
│   ├── crud/           # CRUD operations on models
│   ├── schemas/        # Request/response schemas for API endpoints
│   ├── endpoints/      # API endpoints (organized by version)
│   └── middlewares/    # Custom middleware such as auth and error handling
├── docs/               # OpenAPI documentation generated by Swaggo
├── Dockerfile          # Dockerfile for building the application
├── docker-compose.yml  # Docker Compose configuration for multi-service setup
├── go.mod              # Go module file
└── README.md           # Project README file
```

🛠️ Getting Started

Prerequisites
Docker: Make sure you have Docker and Docker Compose installed on your machine.
Go: Go should be installed if you want to run the application locally without Docker (optional).

🚀 Cloning the Repository
To clone the repository, use the following command:

```bash
git clone https://github.com/omar-aldesi/go-ecommerce.git
cd ecommerce
```

🐳 Running the Application with Docker
The easiest way to run the project is using Docker and Docker Compose. This will automatically set up the application along with a PostgreSQL database.

Build and Run the Docker Containers:

```bash
docker-compose up --build
```

Access the Application:

The API will be available at http://localhost:8080.

Access Swagger UI:

The OpenAPI documentation will be accessible at http://localhost:8080/swagger/index.html.

⚙️ Running the Application Locally
If you prefer to run the application locally without Docker:

Set Up PostgreSQL:

Install PostgreSQL and create a database for the project.
Update the database connection details in the environment variables.
Install Dependencies:

```bash
go mod download
Run the Application:
```

```bash
go run main.go
```

Access the Application:

The API will be running at http://localhost:8080.

🛡️ Authentication
This API uses JWT tokens for authentication. Clients must include a valid JWT in the Authorization header for secured endpoints. The format should be:

```http
Authorization: Bearer <your-jwt-token>
```

🔑 JWT Token Management
Login: Clients can obtain a JWT by providing valid credentials via the /api/v1/auth/login endpoint.
Refresh Token: The application supports token refreshing and blacklisting, ensuring a secure token lifecycle.

🧩 API Documentation
The API documentation is generated using OpenAPI (Swagger) and is accessible at:

```bash
http://localhost:8080/swagger/index.html
```

This documentation provides detailed information about each endpoint, including request parameters, response formats, and authentication requirements.

⚠️ Error Handling
The application has a custom error handling system that provides detailed and user-friendly error messages. Here are some examples:

Validation Error Handling
When a request fails validation, the system returns an error response with details about each invalid field:

```go
// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

// HandleValidationErrors processes and returns validation errors
func HandleValidationErrors(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	var je *json.UnmarshalTypeError

	switch {
	case errors.As(err, &ve):
		errs := make(map[string]interface{})
		for _, e := range ve {
			errs[e.Field()] = formatErrorMessage(e)
		}
		c.JSON(400, ErrorResponse{Errors: errs})
	case errors.As(err, &je):
		errs := map[string]interface{}{
			je.Field: fmt.Sprintf("Invalid value type. Expected %s", je.Type.String()),
		}
		c.JSON(400, ErrorResponse{Errors: errs})
	default:
		c.JSON(400, ErrorResponse{Errors: map[string]interface{}{"general": err.Error()}})
	}
}

// formatErrorMessage formats a single validation error message
func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("This field must be at least %s characters long", e.Param())
		}
		return fmt.Sprintf("This field must be at least %s", e.Param())
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("This field must be at most %s characters long", e.Param())
		}
		return fmt.Sprintf("This field must be at most %s", e.Param())
	case "e164":
		return "Invalid phone number format"
	case "oneof":
		return fmt.Sprintf("This field must be one of: %s", strings.Replace(e.Param(), " ", ", ", -1))
	case "len":
		return fmt.Sprintf("This field must be exactly %s characters long", e.Param())
	case "numeric":
		return "This field must contain only numeric characters"
	case "alphanum":
		return "This field must contain only alphanumeric characters"
	default:
		return fmt.Sprintf("Invalid value for %s", e.Field())
	}
}
```

Custom HTTP Error Handling
You can use the HTTPError struct to return custom HTTP error responses:

```go
// HTTPError is a custom error type for HTTP errors
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func CustomErrorResponse(c *gin.Context, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode, gin.H{"error": httpErr.Message})
	}
	log.Println("Error --> ", httpErr.StatusCode, err)
}
```
This allows for a consistent and clear error response structure across the application.

🧪 Testing
You can run unit and integration tests for the application using:

```bash
go test ./...
```

Make sure to configure the test database connection in your environment variables before running tests.

🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request with your changes.

📜 License
This project is licensed under the MIT License.
