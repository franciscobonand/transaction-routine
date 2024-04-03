# Project transaction-routine

This project implements a service that handles customer finantial transactions as detailed below:

- Each customer has an account with their data.  
- For each operation carried out by the customer, a transaction is created and associated with their respective account.  
- Each transaction has a type (cash purchase, installment purchase, withdrawal or payment), a value and a creation date. Purchase and withdrawal type transactions are recorded with a negative value, while payment transactions are recorded with a positive value.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## MakeFile

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make db-up
```

Shutdown DB container
```bash
make db-down
```

Apply DB migrations
```bash
make migrate-up
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```