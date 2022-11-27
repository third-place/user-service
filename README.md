# Otto User Service

What is it?

An open source identity and access management microservice. Written in go, ready to run in AWS Lambda and backend by 
AWS Cognito.

## Running the server
To run the server, follow these steps:

Create an .env file:
```.env
COGNITO_CLIENT_ID=
COGNITO_CLIENT_SECRET= # if a secret is defined
USER_POOL_ID=
AWS_REGION=

PG_HOST=
PG_USER=
PG_PORT=
PG_PASSWORD=
PG_DBNAME=
```

Start the Postgres database (`-d` runs it as a background process):
```
docker-compose up -d
```

Have AWS keys in a default credential chain.

Run the migrations:
```
go run migrations/*.go
```

Run the server locally:
```
go run main.go
```

Run the tests:
```
go test ./internal/...
```

## Development

Generating models:
```
./bin/swagger-generate-models
```

## Sign Up Flow

![Sign up flow](https://github.com/danielmunro/otto-user-service/blob/main/ref/sign-up.png?raw=true)

## Todo

* better error handling
* groups
* related entities (email, password)
* versioned docs
* recruit contributors
* security audit
* more diagrams

