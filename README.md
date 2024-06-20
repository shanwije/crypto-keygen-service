# Crypto Keygen Service

## Overview

Crypto Keygen Service is a backend service that exposes a REST API for generating Bitcoin and Ethereum crypto addresses.
The service deterministically generates and securely persists the addresses along with their corresponding public and
private keys. This design ensures consistent key generation, easy retrieval, and robust security, making it ideal for
applications requiring reliable and efficient key management.

## Features

- Generate deterministic Bitcoin and Ethereum addresses. ( Designed for extensibility to support more networks in the
  future.)
- Persists the generated addresses and keys ( private key is encrypted).
- Return the generated / persisted address, public key, and private key.
- Includes unit tests to ensure correctness.

## Project Structure

```
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── crypto-keygen-service
│       └── main.go
├── crypto-keygen-service
├── docker-compose.yml
├──
go
.mod
├──
go
.sum
├── internal
│   ├── handler
│   │   └── keygen.go
│   ├── repositories
│   │   └── keygen.go
│   ├── services
│   │   ├── keygen.go
│   │   └── keygen_test.go
│   └── util
│       ├── currency_network_factory
│       │   ├── factory.go
│       │   ├── factory_test.go
│       │   ├── generator.go
│       │   └── generators
│       │       ├── bitcoin
│       │       │   ├── bitcoin.go
│       │       │   └── bitcoin_test.go
│       │       └── ethereum
│       │           ├── ethereum.go
│       │           └── ethereum_test.go
│       ├── encryption
│       │   ├── encryption.go
│       │   └── encryption_test.go
│       └── errors
│           └── errors.go
└── mongo.conf

14 directories, 23 files

```

## Getting Started

### Prerequisites

- Go (1.22.4)
- Docker
- Docker Compose

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/shanwije/crypto-keygen-service
   cd crypto-keygen-service

2. Create a .env file in the root directory with the following content ( refer to .env.example, the given encryption
   key, and seed
   are sample keys. Do not use the same in prod):

   ```bash
    MONGODB_URI=mongodb://localhost:27017
    SERVER_PORT=8080
    DB_NAME=crypto-keygen-service
    DB_COLLECTION=crypto-wallet-service
    GIN_MODE=debug
    MASTER_SEED=6A9D8F4B3C7E1F9A2B8C5D4E7F3A1B2C3D4E5F6A7B8C9D0E1F2A3B4C5D6E7F8
    ENCRYPTION_KEY=4GRrhM8ClnrSmCrDvyFzPKdkJF9NcRkKwxlmIrsYhx0=
   ```
   
    Note: to remove the following warning, set GIN_MODE=release in the .env file
    ```
    crypto-keygen-service-1  | [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    crypto-keygen-service-1  |  - using env:        export GIN_MODE=release
    crypto-keygen-service-1  |  - using code:       gin.SetMode(gin.ReleaseMode)
    ```

3. Build and Run

   Using Docker
   Build and run the application:
    ```bash
    docker-compose up --build
    ```

   Run tests ( require mongo running):
    ```bash
    make test
    ```
   Run the application:
    ```bash
    make run
    ```

## API Endpoints

## Generate / Get Keys and Address

- **URL:** `/keygen/:userId/:network`
- **Method:** `GET`
- **URL Parameters:**
    - `userId` (int): User ID
    - `network` (string): Network type ( bitcoin or ethereum )

### API Responses

#### Success Response

- **Code:** 200
- **Content:**
  ```json
  {
    "address": "generated_address",
    "public_key": "generated_public_key",
    "private_key": "generated_private_key"
  }

#### Error Responses

- **Code:** 400 Bad Request

    - **Content:**
      ```json
      {
        "error": "userId must be a positive integer"
      }
      ```
        - **Possible reasons:**
            - `userId` is not a positive integer.

    - **Content:**
      ```json
      {
        "error": "Network is required"
      }
      ```
        - **Possible reasons:**
            - `network` parameter is missing.

    - **Content:**
      ```json
      {
        "error": "Validation error: [specific error details]"
      }
      ```
        - **Possible reasons:**
            - Specific validation errors related to `userId` and `network` parameters.


- **Code:** 500 Internal Server Error

    - **Content:**
      ```json
      {
        "error": "Internal server error"
      }
      ```
        - **Possible reasons:**
            - Unexpected errors during key generation or database operations.
            - Issues with encrypting/decrypting private keys.

## Health Check

- **URL:** `/health`
- **Method:** `GET`
- **Success Response:**
    - **Code:** 200
    - **Content:** `{"status": "ok"}`
- **Error Response:**
- **Code:** 503 Service Unavailable
    - **Content:** `{"status": "error", "error": [error_message]}`
    - **Possible reasons:**
        - Database connection issues.
        - Service not running.

## Project Components

### Handler

Handles incoming HTTP requests and calls the service layer to process the request.

### Services

Contains the business logic for generating keys and addresses.

### Repositories

Interacts with the MongoDB database to store and retrieve keys.

### Utilities

- **Encryption:** Provides encryption and decryption functionalities.
- **Currency Network Factory:** Uses the factory pattern to generate keys for different networks.
- **Errors:** Defines custom error types for the application.

### Makefile

Provides commands to build, run, test, and clean the project.

## Key Persistence Rationale

### Justification for Persisting Keys

Persisting keys in a database offers several advantages despite the deterministic generation of keys.

#### Advantages

1. **Ease of Access**:
    - Quick retrieval of keys, reducing latency and improving performance.

2. **State Management**:
    - Simplifies the design of stateless services, supporting horizontal scaling.

3. **Audit and Compliance**:
    - Provides an audit trail for key creation, access, and usage, aiding compliance.

4. **Backup and Recovery**:
    - Ensures continuity of service through robust backup and recovery mechanisms.

5. **Reduced Computation**:
    - Avoids the computational expense of regenerating keys for each request.

#### Trade-Offs and Risks

1. Security Risks:
    - **Data Breaches**
    - **Encryption Overhead**

2. Operational Complexity:
    - Requires robust security measures, including encryption at rest and in transit, access controls, and regular
      audits.
    - Encryption and decryption operations can introduce performance overhead.

3. Compliance Requirements ( the biggest headache :) ):
    - Handling and storing sensitive data may require compliance with regulations, which can be complex and costly.

### Security Measures

1. **Encryption**:
    - Encrypt private keys before storage using strong algorithms.
    - Manage encryption keys securely via environment variables or a Key Management Service (KMS).

- ( Below only applicable to production environment )

2. **Access Control**:
    - Implement strict access controls to limit database access.
    - Enforce least privilege access with roles and permissions.

3. **Audit Logging**:
    - Enable audit logging to track access and modifications, with regular reviews for suspicious activity.

4. **Secure Environment**:
    - Protect database and application servers with firewalls, network segmentation, and other security measures.

By balancing these advantages with the trade-offs and implementing robust security measures, persisting keys in a
database can effectively manage key access and ensure service reliability.

### Deterministic Key Generation

Keys are generated deterministically based on the user ID and a master seed, ensuring consistent key pairs for the same
user ID and network.

#### How It was impelemented

1. **Master Seed**: Securely stored and used to generate user-specific seeds.
2. **User-Specific Seed**: Generated using HMAC-SHA256 with the master seed and user ID.
3. **Key Pair Generation**: The user-specific seed generates deterministic key pairs for Bitcoin and Ethereum.

### Why I've persisted the deterministic Keys

1. **Consistency and Performance**:
    - Avoids the computational cost of regenerating keys, improving performance under load.

2. **Enhanced User Experience**

3. **Redundancy and Recovery**

Persisting deterministic keys balances performance, ease of access, and operational simplicity with robust security,
ensuring a reliable and efficient key management service.
