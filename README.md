<p align="center">
  <img src="https://github.com/user-attachments/assets/d5bd1aae-85c7-4ac3-a8e1-ab843d738c8b" alt="VHS Club Logo" width="400"/>
</p>

<p align="center">
  <strong>VHS movie rental service</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#demo">Demo</a> •
  <a href="#installation">Installation</a> •
  <a href="#ci">CI</a>
</p>

---

<img alt="vhs-club header" src="docs/assets/vhs-header2-1080p.gif" />

## Overview

VHS Club is a Linux full-stack web application for a fictional VHS movie rental service. Users can browse a catalog of classic VHS movies and rent their favorite titles. The application comes with pre-created users and an admin account loaded from SQL seed scripts.

## Features

- **User Authentication**: Secure login with JWT-based authentication
- **VHS Catalog**: Browse a collection of VHS movies with details (title, director, genre)
- **Rental System**: Rent available tapes and track rental history
- **Role-Based Access**: Admin and user roles with different permissions
- **React Frontend**: User interface built with React for a smooth browsing and rental experience
- **RESTful API**: Clean API design following REST principles
- **Dockerized Deployment**: Complete containerization with Docker Compose for easy setup and deployment
- **Database Seeding**: SQL scripts provide initial data including admin account, sample users, and VHS tapes

## Tech Stack

### Backend

| Technology | Description |
|------------|-------------|
| **Go** | Backend language (v1.25) |
| **Gin Gonic** | High-performance HTTP web framework |
| **PostgreSQL** | Relational database for data persistence |
| **JWT** | Authentication tokens |
| **Argon2id** | Secure password hashing |
| **SQLC** | Type-safe SQL code generator |

### Frontend

| Technology | Description |
|------------|-------------|
| **React** | UI library (v19) |
| **Vite** | Build tool and dev server |
| **React Router** | Client-side routing |
| **Vitest** | Unit testing framework |
| **Nginx** | Production web server |

### DevOps & Tools

| Technology | Description |
|------------|-------------|
| **Docker** | Containerization |
| **Docker Compose** | Multi-container orchestration |
| **GitHub Actions** | CI/CD pipeline |
| **Nginx** | Reverse proxy and static file serving |

## Demo

### Screenshots

<p align="center">
  <img src="https://github.com/user-attachments/assets/08cd0299-bfd7-4f71-9be9-ac845c8ac387" alt="Home Page" width="410"/>
</p>
<p align="center">
  <img src="https://github.com/user-attachments/assets/25700c50-fed4-498e-9e48-d865b7374384" alt="Catalog Page" width="410"/>
</p>

### Video Demo

<video src="https://github.com/user-attachments/assets/62a0c20e-715d-46c7-ac17-965eb01ff301" width="400" controls></video>

## Installation

### Prerequisites

- **Docker** (v20.10+) and **Docker Compose** (v2.0+) - for containerized setup
- **Go** (v1.25+) - for local backend development
- **Node.js** (v24+) and **npm** - for local frontend development
- **PostgreSQL** (v16+) - for local database

### Option 1: Docker Compose (Recommended)

The easiest way to run the entire application stack is using Docker Compose. This will spin up all three services (database, backend, and frontend) with a single command.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/rigofekete/vhs-club.git
   cd vhs-club
   ```

2. **Set up environment variables:**

   Create a `.env` file in the project root with a JWT secret:

   ```bash
   # Generate a secure JWT secret
   JWT_SECRET=$(openssl rand -base64 32)

   # Create the .env file
   echo "JWT_SECRET=$JWT_SECRET" > .env
   ```

   The database connection URL is handled internally by Docker Compose, but the JWT secret must be provided by you.

3. **Build and run with Docker Compose:**

   ```bash
   docker compose up -d --build
   ```

   This command will:
   - Build the database image with initial schema and seed data (includes admin account and sample users)
   - Build the Go backend API
   - Build the React frontend with Nginx
   - Start all services with proper networking

   <br>

   > **Important:** The database seeding happens automatically on first startup. The `Dockerfile.db` copies `sql/seed.sql` into the PostgreSQL initialization directory, which PostgreSQL executes when the container starts for the first time.

4. **Access the application:**

   Open your browser and navigate to: **<http://localhost>**

   | Service | URL | Description |
   |---------|-----|-------------|
   | Frontend | <http://localhost> | Main web application |
   | Backend API | <http://localhost:8080> | REST API endpoints |
   | Health Check | <http://localhost:8080/health> | API health status |

5. **Stop the services:**

   ```bash
   docker compose down
   ```

   To remove all data (including the database volume):

   ```bash
   docker compose down -v
   ```

### Option 2: Local Development Setup (Without Docker)

If you prefer to run the services directly on your machine for development purposes, follow these steps.

#### Step 1: Database Setup

1. **Install PostgreSQL** (v16+):

   On Arch Linux:

   ```bash
   sudo pacman -S postgresql
   sudo systemctl enable postgresql
   sudo systemctl start postgresql
   ```

   > **Note:** These instructions have been tested on Arch Linux. For other distributions, consult your package manager's documentation for PostgreSQL installation.

2. **Run the seed script:**

   The `sql/seed.sql` file handles everything: it creates the database, schema, tables, and populates them with initial data (admin account, sample users, and VHS tapes).

   ```bash
   # Run as postgres superuser - this creates the database and everything else
   sudo -u postgres psql -f sql/seed.sql
   ```

   > **Important:** The seed script creates the `vhs_club` database, all tables, and initial data in one step. After seeding, you can log in with the default credentials listed in the [Database Seeding](#database-seeding) section.

#### Step 2: Backend Setup

1. **Install Go** (v1.25+):

   - Download from [https://go.dev/dl/](https://go.dev/dl/)
   - Or use your package manager

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Configure the database connection:**

   Create a `.env` file in the project root with your PostgreSQL credentials:

   ```bash
   # Replace with your actual PostgreSQL credentials
   DB_URL=postgres://your_username:your_password@localhost:5432/vhs_club?sslmode=disable
   JWT_SECRET=your-secret-key-here
   ```

   For example, if your PostgreSQL user is `postgres` with password `postgres`:

   ```bash
   DB_URL=postgres://postgres:postgres@localhost:5432/vhs_club?sslmode=disable
   JWT_SECRET=your-secret-key-here
   ```

   **Finding your PostgreSQL credentials:**

   - If you installed PostgreSQL locally, the default superuser is usually `postgres` with a password you set during installation
   - Check your PostgreSQL configuration in `pg_hba.conf` for authentication methods
   - If you seeded using `sudo -u postgres psql -f sql/seed.sql`, the database was created by the `postgres` superuser

   Or export them directly in your terminal:

   ```bash
   export DB_URL="postgres://your_username:your_password@localhost:5432/vhs_club?sslmode=disable"
   export JWT_SECRET="your-secret-key-here"
   ```

4. **Run the backend:**

   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`

5. **Run tests:**

   ```bash
   go test ./...
   ```

#### Step 3: Frontend Setup

1. **Install Node.js** (v24+):

   - Download from [https://nodejs.org/](https://nodejs.org/)
   - Or use nvm: `nvm install 24`

2. **Navigate to the frontend directory:**

   ```bash
   cd react-app
   ```

3. **Install dependencies:**

   ```bash
   npm install
   ```

4. **Start the development server:**

   ```bash
   npm run dev
   ```

   The frontend will be available at `http://localhost:5173` (default Vite port)

5. **Build for production:**

   ```bash
   npm run build
   ```

6. **Run tests:**

   ```bash
   npx vitest run
   ```

### Development URLs Summary

| Service | Local Development | Docker |
|---------|-------------------|--------|
| Frontend | <http://localhost:5173> | <http://localhost> |
| Backend API | <http://localhost:8080> | <http://localhost:8080> |
| Database | localhost:5432 | localhost:5432 |

## API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users/login` | Authenticate and receive JWT token |

> **Note:** User accounts are pre-created via SQL seed scripts. There is no public registration endpoint. However, admin users can create new accounts individually or in batches through the [User Management Endpoints](#user-management-endpoints-admin). See the [Database Seeding](#database-seeding) section for default credentials.

### Tapes Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tapes` | List all available tapes (public) |
| GET | `/api/tapes/:id` | Get a specific tape by ID (public) |
| POST | `/api/tapes` | Create a new tape (admin only) |
| POST | `/api/tapes/batch` | Create multiple tapes (admin only) |
| PATCH | `/api/tapes/:id` | Update a tape (admin only) |
| DELETE | `/api/tapes/:id` | Delete a tape (admin only) |
| DELETE | `/api/tapes` | Delete all tapes (admin only) |

### Rentals Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rentals` | List all active rentals (public) |
| POST | `/api/rentals/:id` | Create a new rental (authenticated users) |
| PATCH | `/api/rentals/:id` | Return a rented tape (authenticated users) |
| DELETE | `/api/rentals` | Delete all rentals (admin only) |

### User Management Endpoints (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create a new user (admin only) |
| POST | `/api/users/batch` | Create multiple users (admin only) |
| GET | `/api/users` | List all users (admin only) |
| GET | `/api/users/:id` | Get user by ID (admin only) |
| DELETE | `/api/users` | Delete all users (admin only) |

## Database Seeding

The application does not have a public user registration endpoint. Instead, users are pre-created via SQL seed scripts. Additional users can be created by an admin via the [User Management Endpoints](#user-management-endpoints-admin) (individually or in batches). The **same `sql/seed.sql` file** is used for both Docker and local development, ensuring consistency across environments. The database is automatically populated with sample data when using Docker Compose, or you can manually apply the seed script for local development.

### Default Credentials

The following accounts are available after seeding:

| Username | Email | Password | Role |
|----------|-------|----------|------|
| **Admin** | `admin@vhs-club.hu` | `12345678` | admin |
| **ArthurCClarke** | `thesentinel@space.odissey` | `12345678` | user |
| **MilesDavis** | `grumpy.genius@cool.com` | `12345678` | user |

### Sample Tapes

The seed script also populates the catalog with classic VHS movies:

| Title | Director | Genre | Quantity |
|-------|----------|-------|----------|
| Amarcord | Federico Fellini | Drama | 1 |
| Taxi Driver | Martin Scorsese | Thriller | 2 |
| Back to the Future | Robert Zemeckis | Adventure | 5 |
| Alien | Ridley Scott | Horror | 10 |
| A torinói ló | Béla Tarr | Drama | 3 |
| Batman | Tim Burton | Action | 4 |
| Fitzcarraldo | Werner Herzog | Drama | 11 |

### Seeding with Docker Compose

When using Docker Compose, the database is automatically seeded on first startup using the same `sql/seed.sql` file used for local development. The `Dockerfile.db` copies the seed script into the PostgreSQL initialization directory, which PostgreSQL executes when the container starts for the first time.

```bash
# Start all services - seeding happens automatically
docker compose up -d --build
```

> **Note:** The database is only seeded on the **first** container startup. If you need to re-seed, you must remove the volume: `docker compose down -v` and then start again.

### Seeding Manually (Local Development)

For local development without Docker, run the **same seed script** used by Docker Compose:

```bash
# This creates the database, schema, tables, and all initial data in one step
sudo -u postgres psql -f sql/seed.sql
```

The `sql/seed.sql` file handles everything:

- Creates the `vhs_club` database
- Creates all tables (users, tapes, rentals)
- Inserts the admin account and sample users
- Inserts the VHS tape catalog

> **Note:** The seed script uses `CREATE DATABASE` and `\c` (connect) commands, so it must be run as a superuser (e.g., `postgres` user) and not as a regular database user.

## CI

### GitHub Actions Pipeline

The project includes a CI pipeline defined in `.github/workflows/ci.yml` that runs on every push to the `master` branch.

### Pipeline Stages

```mermaid
graph LR
    A[Push/PR] --> B[Test Backend]
    A --> C[Test Frontend]
    B --> D[Build Images]
    C --> D
    D --> E[Smoke Test]
```

1. **Test Backend**: Runs all Go unit tests
2. **Test Frontend**: Runs all Vitest tests for the React application
3. **Build**: Builds Docker images for all three services (backend, frontend, database)
4. **Smoke Test**: Deploys the full stack and verifies the application is accessible

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_URL` | PostgreSQL connection string | Yes | - |
| `JWT_SECRET` | Secret key for JWT signing | Yes | - |

### Generating a JWT Secret

For production deployments, generate a secure random secret:

```bash
# Generate a 32-byte base64-encoded secret
openssl rand -base64 32
```

## Development

### Project Architecture

The backend follows a layered architecture pattern:

```
┌─────────────────┐
│   Handler       │  ← HTTP request handlers (Controllers)
│   (Gin)         │
├─────────────────┤
│   Service       │  ← Business logic
│                 │
├─────────────────┤
│   Repository    │  ← Data access layer
│                 │
├─────────────────┤
│   Database      │  ← PostgreSQL
│   (PostgreSQL)  │
└─────────────────┘
```

### Backend Development

```bash
# Run the server with hot reload (requires air or similar)
go run main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate SQLC code (after modifying SQL queries)
sqlc generate
```

### Database Schema & Code Generation

The backend uses [SQLC](https://sqlc.dev/) to generate type-safe Go code from SQL queries. This ensures compile-time type safety and eliminates the need for manual struct mapping.

#### Schema Location

The database schema and queries are organized as follows:

| Path | Description |
|------|-------------|
| `sql/schema/` | Database schema migration files (001_tapes.sql, 002_users.sql, 003_rentals.sql) |
| `sql/queries/` | SQL queries used by the application (users.sql, tapes.sql, rentals.sql) |
| `sql/seed.sql` | Complete seed script that creates the database, schema, tables, and initial data |
| `sqlc.yaml` | SQLC configuration file |

> **Important:** SQLC reads all files in the `sql/schema/` directory to understand the database structure when generating Go types. The `sqlc.yaml` file points SQLC to these paths.

#### Two Approaches to Database Setup

You have two options for working with the database:

**Option 1: Use the Ready-Made Seed (Recommended for Quick Start)**

The easiest way is to use the provided `sql/seed.sql` script, which creates the database and populates it with initial data:

```bash
# Run the complete seed script (creates database, schema, tables, and sample data)
sudo -u postgres psql -f sql/seed.sql
```

> **Note:** This approach uses the pre-generated Go types in `internal/database/`. No additional steps are needed — the generated code is already committed to the repository and ready to use.

**Option 2: Modify the Schema and Regenerate**

If you want to customize the database structure:

1. **Edit the schema files** in `sql/schema/` (001_tapes.sql, 002_users.sql, 003_rentals.sql) to add or modify tables, columns, or constraints

2. **Edit or add queries** in the `sql/queries/` directory (e.g., `users.sql`, `tapes.sql`)

3. **Regenerate the Go code**:

   Run the following commands from the **project root** (where `sqlc.yaml` is located):

   ```bash
   # Install SQLC if you haven't already
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

   # Generate type-safe Go code from SQL queries
   sqlc generate
   ```

   This command reads `sqlc.yaml` and generates Go structs and methods in `internal/database/`.

4. **Apply the updated schema** to your database by running the modified schema files or re-seeding

> **Important:** When modifying the schema, you **must** run `sqlc generate` afterwards. The Go code in `internal/database/` is auto-generated from the SQL files, and any changes to the schema will not be reflected in the Go types until you regenerate.

### Frontend Development

```bash
cd react-app

# Start development server
npm run dev

# Run tests
npx vitest run

# Run tests in watch mode
npx vitest

# Build for production
npm run build

# Preview production build
npm run preview
```

## Testing

### Unit Tests

```bash
# Run all Go tests with verbose output
go test -v ./...

# Run tests with coverage report
go test -cover ./...
```

### Frontend Tests

```bash
cd react-app

# Run Vitest tests
npx vitest run

# Run with coverage
npx vitest run --coverage

# Run in UI mode
npx vitest --ui
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
