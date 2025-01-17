# event-management

# Event Management API

An API designed to manage events, supporting features like event creation, user authentication, ticket generation, and payment processing. This project follows modern software development practices and uses robust technologies to ensure scalability and reliability.

---

## Functional Requirements

### 1. Event Creation
- **Endpoint to create an event** with the following fields:
  - Event Title
  - Description
  - Location
  - Start and End Time
  - Price (if paid)
  - Number of tickets available

### 2. Event Management
- **Update and delete events** created by the user.
- **List events** belonging to the authenticated user.

### 3. Payments (Optional in MVP)
- Process payments for paid events using **Stripe**.

### 4. Ticket Generation
- Generate unique tickets for participants.
- Store generated tickets for future validation.

### 5. Ticket Validation
- Validate tickets using QR code or unique codes.

---

## Non-Functional Requirements

- RESTful API design.
- Adequate performance to handle **1000 concurrent users**.
- Comprehensive API documentation.

---

## Technologies Used

- **Backend Framework**: [Gin](https://gin-gonic.com/)
- **Database**: PostgreSQL
- **Authentication**: OAuth2 (Google, GitHub) via [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- **Payments**: [Stripe](https://stripe.com/)
- **API Documentation**: Swagger ([Swaggo](https://github.com/swaggo/swag))
- **Testing**:
  - Unit and Integration: [Testify](https://github.com/stretchr/testify)
  - Load Testing: [k6](https://k6.io/) or [Artillery](https://www.artillery.io/)
- **Containerization**: Docker
- **Optional**:
  - Monitoring: Grafana + Prometheus
  - Orchestration: Kubernetes

---

## Roadmap

### **Phase 1: Fundamentals**
1. **Initial Setup**
   - [x] Configure Gin.
   - [x] Set up basic project structure.
   - [x] Connect to PostgreSQL and configure migrations using `golang-migrate`.
2. **Documentation**
   - [x] Set up Swagger for API documentation.
   - [x] Draw the architecture with tools like Excalidraw.
   - [x] Draw or write the Use Cases .
3. **Basic Endpoints**
   - [x] Create Event.
   - [x] List Events.
   - - [x] Create a helper for pagination
   - [x] Create a BaseResponse and Metadata about pagination.
4. **Containerization and Deployment**
   - [x] Configure Docker for development and production environments.
   - [x] Create CD rotine to first deploy.
   - [x] Make a deployment.
5. **Authentication** (Needs first deploy and a domain to simplify OAuth2 google configuration to test...)
   - [x] Configure in GCP.
   - [x] Implement OAuth2 login with Google using OpenID for authentication.
   - [x] Return the JWT token if the authentication is successful.
   - [x] Refactor: Change to Struct-Based Controller,Service and Repository.
   - [ ] Refresh token endpoint.
   - [ ] Private routes.

### **Phase 2: Intermediate Features**
1. **Event Management**
   - [ ] Update and delete events.
   - [ ] Add pagination for event listings.
   - [ ] Protect endpoints with Private endpoints.
2. **Payments**
   - [ ] Integrate Stripe for processing payments (Checkout Sessions).
3. **Ticket Handling**
   - [ ] Generate unique tickets upon successful payment.
   - [ ] Store ticket details in the database.

### **Phase 3: Validation and Deployment**
1. **Ticket Validation**
   - [ ] Implement QR code simulation or unique code validation.
2. **Testing**
   - [ ] Write unit and integration tests.
   - [ ] Perform load testing with k6 or Artillery.
3. **CI/CD**
   - [ ] Create Github action CI.
   - [ ] Create a workflow to automate releases (using semantic commit).
4. **Optional**
   - [ ] Create Kubernetes deployments and services.
   - [ ] Integrate Grafana and Prometheus for monitoring.

---

## Architecture Overview

![image](https://github.com/user-attachments/assets/0d139d49-88d8-4d7d-bff6-6ffb6dac0aa7)


### Components
- **Controllers**: Handle HTTP requests and responses.
- **Services**: Contain business logic.
- **Repositories**: Handle database interactions.
- **Models**: Define data structures.


### Development Commands

- `./scripts/migrate.sh up`
- `./scripts/migrate.sh down`
- `./scripts/migrate.sh version`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any suggestions or improvements.

