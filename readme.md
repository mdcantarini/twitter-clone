# Twitter Clone - Microservices Architecture

A distributed Twitter clone built with Go, implementing a microservices architecture with event-driven communication.

## Project Architecture

This project implements a simplified Twitter clone using a microservices architecture with the following components:

### Microservices

1. **User API** (Port 8081)
   - User registration and profile management
   - SQLite database for user data
   - Endpoints: Create user, Get user by ID

2. **Tweet API** (Port 8082)
   - Tweet creation and retrieval
   - Cassandra database for scalable tweet storage
   - Publishes tweet events to Kafka for fan-out
   - Endpoints: Create tweet, Get tweet by ID

3. **Follow API** (Port 8083)
   - Manage follow relationships between users
   - SQLite database for follow data
   - Endpoints: Follow user, Get follower IDs

4. **Feed API** (Port 8084)
   - User timeline/feed generation
   - Consumes tweet events from Kafka
   - Cassandra database for pre-computed user timelines
   - Implements fan-out on write pattern
   - Endpoints: Get user feed

### Infrastructure Components

- **Cassandra**: NoSQL database for tweets and user timelines
- **SQLite**: Relational database for users and follow relationships
- **Kafka**: Message broker for asynchronous tweet event processing
- **Zookeeper**: Kafka coordination

### Architecture Patterns

- **Microservices**: Each service is independently deployable with its own database
- **Event-Driven**: Tweet creation triggers asynchronous fan-out via Kafka
- **Fan-out on Write**: Pre-compute user timelines for fast feed retrieval
- **Repository Pattern**: Clean separation between business logic and data access

## How to Run

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

### Starting the Application

1. Clone the repository
2. Create environment files for each service (or copy from examples):
   ```bash
   # Create env files from the example
   cp .env.example .env
   ```
3. Run the entire stack:
   ```bash
   docker-compose up -d
   ```
4. Wait for all services to be healthy (check with `docker-compose ps`)

### Stopping the Application

```bash
docker-compose down
```

To also remove volumes:
```bash
docker-compose down -v
```

## API Endpoints

### User API (http://localhost:8081)

**Create User**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "johndoe",
  "display_name": "John Doe"
}
```

**Get User**
```bash
GET /api/v1/users/{id}
```

### Tweet API (http://localhost:8082)

**Create Tweet**
```bash
POST /api/v1/tweets
Content-Type: application/json

{
  "user_id": 1,
  "content": "Hello Twitter!"
}
```

**Get Tweet**
```bash
GET /api/v1/tweet/{tweet_id}
```

### Follow API (http://localhost:8083)

**Follow User**
```bash
POST /api/v1/follow
Content-Type: application/json

{
  "follower_id": 1,
  "followed_id": 2
}
```

**Get Follower IDs**
```bash
GET /api/v1/users/{user_id}/follower_ids
```

### Feed API (http://localhost:8084)

**Get User Feed**
```bash
GET /api/v1/feed/{user_id}
```

## Areas of Improvement Detected

### 1. Database Technology
- **Issue**: Using SQLite for user and follow data limits scalability
- **Recommendation**: Replace SQLite with PostgreSQL for better concurrency and production readiness
- **Location**: `docker-compose.yml:37` comment already identifies this

### 2. Error Handling & Reliability
- **Tweet Creation**: Tweet is saved but fan-out can fail silently (`internal/tweet/service.go:82`)
- **Message Processing**: Lost messages if processing fails (`internal/feed/service.go:86`)
- **Recommendation**: Implement retry mechanisms, dead letter queues, and transactional outbox pattern

### 3. Caching
- **Missing cache layer** for hot tweets (`internal/tweet/service.go:99`)
- **Recommendation**: Add Redis for caching frequently accessed tweets and user profiles

### 4. Pagination
- **Hard limit** of 50 entries in user feed (`internal/feed/service.go:44`)
- **Recommendation**: Implement proper pagination with cursor-based navigation

### 5. Security & Authentication
- **No authentication** mechanism implemented
- **No rate limiting**
- **Recommendation**: Add JWT-based auth, API rate limiting, and input sanitization

### 6. Testing
- **Limited test coverage** - only service-level tests exist
- **No integration tests** for the full system
- **Recommendation**: Add integration tests, and increase unit test coverage

### 7. Observability
- **No metrics or tracing** implemented
- **Limited logging**
- **Recommendation**: Add Prometheus metrics, and structured logging

### 8. API Design
- **Inconsistent routing**: `/tweet/:tweet_id` vs `/tweets` (singular vs plural)
- **No API versioning** strategy beyond `/api/v1`
- **Recommendation**: Standardize REST conventions and implement proper API versioning

### 9. Config Management
- **Single Env config** for all the services
- **Recommendation**: Implement configs per service