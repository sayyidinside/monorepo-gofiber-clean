# Monorepo GoFiber Clean

Monorepo GoFiber Clean is a Go backend starter following Clean Architecture principles, adapted for a monorepo layout.
The repo is split into services/ (each service is a standalone module) and a shared/ module for reusable domain, infra and helpers. It includes tooling for local dev (Air hot-reload), DI/bootstrapping, migrations, and a recommended workflow for CI/CD.

## Features / Technologies Used

- **GoFiber**: Web framework for building fast and scalable APIs.
- **GORM**: Object-Relational Mapper (ORM) for database operations, utilizing GORM datatypes.
- **Redis**: In-memory key–value database, used as a distributed cache and message broker, with optional durability.
- **Air**: Live reload for Go applications during development.
- **Zap**: Fast, structured logging.
- **Validator V10**: Validation of incoming data for requests.
- **Gocron**: Job scheduling for recurring tasks.
- **Bluemonday**: HTML sanitizer for handling user-generated content securely.
- **Viper**: For configuration management, with support for environment variables and multiple file formats (YAML, JSON, etc.), including auto-reloading of configuration files.
- **Clean Architecture**: A layered approach to structure the codebase for maintainability and scalability.
- **Monorepo**: software-development strategy in which the code for a number of projects is stored in the same repository.

## Project Structure

```bash
├── services/                         # Independent services (each is a separate binary)
│    ├── service_a/                   # Example service folder
│    │    ├── cmd/                    # Main application entry points
│    │    │   ├── server/             # HTTP server setup
│    │    │   ├── worker/             # Background worker setup
│    │    │   ├── bootstrap/          # depedency initialization
│    │    ├── domain/                 # Spesific Core business logic and domain-specific concerns
│    │    │   ├── entity/             # Defines the core business entities (user, role, permission, etc)
│    │    │   ├── repository/         # Defines the interfaces for interacting with data persistence.
│    │    │   └── service/            # Contains the business logic
│    │    ├── interfaces/             # Spesific interface adapters (Delivery layer)
│    │    │   ├── http/               # HTTP delivery (GoFiber routes)
│    │    │   │   ├── auth/           # HTTP handlers for auth-related routes
│    │    │   │   ├── handler/        # General handlers (HTTP request handling logic)
│    │    │   │   ├── permission/     # HTTP handlers for permission-related routes
│    │    │   │   ├── role/           # HTTP handlers for role-related routes
│    │    │   │   ├── routes/         # Route definitions for api
│    │    │   │   │   └── v1/         # Versioned API routes (e.g., v1 API)
│    │    │   │   │       └── users/  # Route related to user management
│    │    │   ├── model/              # Data transfer objects (DTOs) for mapping HTTP <-> domain
├── shared/                           # Code reusable across services
│   ├── domain/                       # Shared Core business logic and domain-specific concerns
│   │   ├── entity/                   # Defines the core business entities (user, role, permission, etc)
│   │   ├── repository/               # Defines the interfaces for interacting with data persistence.
│   │   └── service/                  # Contains the business logic
│   ├── infrastructure/               # Shared infrastructure-specific code (frameworks, DB, etc.)
│   │   ├── config/                   # Configuration files (loading .env variables, app settings)
│   │   ├── database/                 # Database setup and implementations (GORM)
│   │   ├── redis/                    # Redis setup and implementations
│   │   ├── shutdown/                 # Service gracefull shutdown implementation
│   │   ├── logger/                   # Logging setup (zap)
│   │   ├── scheduler/                # Scheduling logic (gocron)
│   │   ├── bootstrap/                # Reusable depedency initialization
│   ├── interfaces/                   # Shared interface adapters (Delivery layer)
│   │   ├── http/                     # HTTP delivery (GoFiber routes)
│   │   │   ├── middleware/           # HTTP middleware (auth, JWT, role-based)
│   │   ├── model/                    # Data transfer objects (DTOs) for mapping HTTP <-> domain
│   ├── pkg/                          # Shared libraries and utilities
│   │   ├── helpers/                  # Generic helper functions (not domain-specific)
│   │   └── utils/                    # Generic utility functions (not domain-specific)
├── tests/                            # Unit and integration tests
├── migration/                        # Database migration and seeding tools
└── .env                              # Environment variables
```

## Installation

1. **Clone the repository:**

```bash
git clone https://github.com/sayyidinside/monorepo-gofiber-clean.git
cd monorepo-gofiber-clean
```

2. **Set up environment variables:**

Create a `.env` file based on `.env.example` and update the configuration as needed.

3. **Sync go workspace:**

```bash
go work sync
```

4. **Install dependencies for each services:**

```bash
cd ./services/service_a
go mod tidy

cd ./services/service_a
go mod tidy
```

5. **Run database migration:**

```bash
cd ./migration
go run main.go
```

6. **Run the application for each services (with live reload):**

```bash
cd ./services/service_a
air

cd ./services/service_b
air
```

## User Management

- **Users**: Manage user accounts.
- **Roles**: Assign different roles to users.
- **Permissions**: Define and assign permissions to roles.

## Auth Middleware

The project includes JWT-based authentication, as well as role-based access control middleware. You can extend the authentication middleware as needed.

## Contributing

Feel free to submit issues or pull requests to improve this project. Make sure to follow the contribution guidelines.

## License

This project is licensed under the MIT License.
