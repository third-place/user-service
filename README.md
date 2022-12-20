# Third Place User Service

This repository represents the identity and access management microservice for [Third place](https://thirdplaceapp.com).

## Tests

Run the tests:
```
go test ./...
```

## OpenAPI

Generating models:
```
./bin/swagger-generate-server.sh
```

## Local Development

1. `docker-compose.yml` is provided, which provisions local versions of kafka, zookeeper, and postgres.

Run `docker-compose up`.

2. Create the kafka topics manually:

```
docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic images

docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic users
```

3. An invite code will be needed later:

```
psql -h localhost -U postgres -p 54321 -c "insert into invites (code) values ('abc-123')"
```

4. Finally, copy the `.env.template` file to `.env` and fill in the secrets.

After the above steps, `go run main.go` should work.

## Sign Up Flow

![Sign up flow](https://github.com/third-place/user-service/blob/main/ref/sign-up.png?raw=true)

## Todo

* better error handling
* groups
* related entities (email, password)
* versioned docs
* recruit contributors
* security audit
* more diagrams

