# E-commerce Microservices
A simple e-commerce platform built with microservices architecture using Go, gRPC, and Docker.

## Architecture
This project implements a microservices-based e-commerce system with the following services:

- User Service - Handles user registration, authentication with JWT, and profile management
- Product Service - Manages product catalog, inventory, and product information
- Order Service - Processes orders, order history, and order status tracking

Each service is completely independent with its own database and communicates via gRPC.

## Tech Stack
- Language: Go 1.24
- Communication: gRPC for inter-service communication
- Containerization: Docker & Docker Compose
- Database: MySQL (per service)
- (To be implemented) Message Queue: Redis (for async operations)

## Future Implementations
- Message Queue: Redis (for async operations)
- Inventory Stock Prediction with Prophet and XGBoost 

## Prerequisites
- Go 1.24 or higher
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)

## Getting Started
1. **Clone the repository**
```bash
git clone <https://github.com/ShernaC/ecommerce.git>
cd ecommerce-microservices
```

2. **Build and run with Docker Compose**
```bash
docker-compose up --build
```

## Environment Variables
Key environment variables for configuration:
```env
# Database
PORT=8080
DB_HOST=host.docker.internal
DB_PORT=3306
USER_DB_DATABASE=ecommerce
PRODUCT_DB_DATABASE=ecommerce_product
ORDER_DB_DATABASE=ecommerce_order

# Service ports
USER_SERVICE_PORT=8080
PRODUCT_SERVICE_PORT=8181
ORDER_SERVICE_PORT=8282

# gRPC ports
USER_GRPC_PORT=50051
PRODUCT_GRPC_PORT=50052
ORDER_GRPC_PORT=50053
```
