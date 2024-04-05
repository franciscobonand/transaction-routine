# Project transaction-routine

This project implements a service that handles customer financial transactions as detailed below:

- Each customer has an account with their data.  
- For each operation carried out by the customer, a transaction is created and associated with their respective account.  
- Each transaction has a type (cash purchase, installment purchase, withdrawal or payment), a value and a creation date. Purchase and withdrawal type transactions are recorded with a negative value, while payment transactions are recorded with a positive value.

## Developed with

- [Golang](https://go.dev/doc/install) (developed with version 1.22)
- [Docker](https://docs.docker.com/engine/install/) and [Docker Compose](https://docs.docker.com/compose/install/)
- Postgres
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Getting Started

### Available endpoints

This application provides several endpoints to manage accounts and transactions.
In the endpoints description below, remember to replace `localhost:8080` with the actual server address and port if different.

#### Health Check

- Endpoint: `/health`
- Method: `GET`
- Description: Checks if the application is running

```bash
curl -X GET http://localhost:8080/health
```

#### Create Account

- Endpoint: `/accounts`
- Method: `POST`
- Description: Creates a new account. The request body should contain the account details in JSON format.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"document_number":"12345678900"}' http://localhost:8080/accounts
```

#### Get Account

- Endpoint: `/accounts/{id}`
- Method: `GET`
- Description: Retrieves the details of an account with the given ID.

```bash
curl -X GET http://localhost:8080/accounts/1
```

#### Get Account Balance

- Endpoint: `/accounts/{id}/balance`
- Method: `GET`
- Description: Retrieves the balance of an account with the given ID.

```bash
curl -X GET http://localhost:8080/accounts/1/balance
```

#### Create Transaction

- Endpoint: `/transactions`
- Method: `POST`
- Description: Creates a new transaction. The request body should contain the transaction details in JSON format.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"account_id":1, "operation_type_id":1, "amount":123.45}' http://localhost:8080/transactions
```

#### Update Transaction

- Endpoint: `/transactions/{id}`
- Method: `PUT`
- Description: Updates a transaction with the given ID. The request body should contain the new transaction details in JSON format.

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"account_id":1, "operation_type_id":1, "amount":123.45}' http://localhost:8080/transactions/1
```

### Running the application

This repo contains a Makefile to manage common tasks such as building, running, and testing the application. Here are the steps to run the application:

1. **Start the database**
The application uses a PostgreSQL database running in a Docker container. Use the `db.up` command to start the database container

```bash
    make db.up
```

2. **Run database migrations**
The application uses database migrations to manage the database schema. Use the `migrate.up` command to apply the migrations.

```bash
make migrate.up
```

3. **Run the application**
Use the `run` command to start the application. This will start the server and it will begin listening for incoming requests.

```bash
make run
```

4. **Live reload (optional)**
If you want the application to automatically rebuild and restart when files change, you can use the `watch` command. This requires the [air](https://github.com/cosmtrek/air) tool to be installed.

```bash
make watch
```

### Testing

#### Unit Test

The unit tests provided in `/tests` focus on validating and testing handlers and the services that contain some additional logic besides sending requests to the repository.  
Tests related to database queries were covered in the **Load Test** section below.

Running the unit tests:

```bash
make test
```

#### Load Test

Load/performance tests were created using the tool [k6](https://grafana.com/docs/k6/latest/).  
Before running this test, be sure to have the database running with the migrations applied, and empty `account` and `transaction` tables.  
It is also expected that the application will be running on `http://localhost:8080`.  

Running the load tests:

```bash
make loadtest
```

## Make - available commands

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Create DB container
```bash
make db.up
```

Shutdown DB container
```bash
make db.down
```

Clear the Docker Volume containing DB data
```bash
make db.clear
```

Reset DB container that is running
```bash
make db.reset
```

Apply DB migrations
```bash
make migrate.up
```

Live reload the application
```bash
make watch
```

Run the test suite
```bash
make test
```

Run load tests
```bash
make loadtest
```

Generate mocks
```bash
make mocks
```

Clean up binary from the last build
```bash
make clean
```
