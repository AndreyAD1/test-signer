# test-signer
The Test signer is a service that accepts a set of answers and questions and signs that the user has finished the " test " at this point in time. The signatures are stored and can later be verified by a different service.

## Getting Started
- Create a PostgreSQL database;
- Install [a migration tool](https://github.com/golang-migrate/migrate);
- Run migrations:
```shell
migrate -database '<db_url>' -path internal/app/infrastructure/migrations up
```
- Install the project dependencies:
```shell 
go mod tidy
```
- Run the server:
```shell 
SIGN_KEY='your secret' go run main.go -u '<db_url>' -s '<a JWT secret>' 
```