# Test backend - Niveau 1

## Exam

1. Objective
   Develop a robust APIs for managing products, categories, suppliers, dynamic filters, optimized
   pagination for continuous loading (infinite scroll).
   The API must be scalable, secure, and provide a smooth interaction with the front-end.
2. Database & Figma
   Database will be provided when talking test.
   Figma Link:
   https://www.figma.com/design/dfTzvBVPK3HdZB5XartT9s/Test---FE?node-id=0-1&t=qvJDo4KCVgkZ7w6x-1
3. API
   Imagine the APIs necessary for the correct display of data as presented on Figma
4. Item Range Management for a Scrollable Board
   Implementation
   - Optimize product loading for a smooth display on a scrollable board.
   - Avoid initial loading of thousands of rows by retrieving products in batches.
   - Dynamic filter and incremental search
   - Create an API that generates a formatted PDF file of product data in the back end.
   - Create an APi to calculate the distance in km between a location (ip) and a city where a
     product produced is located.
5. Date Format Validation (Regex)
   All date fields should follow the YYYY-MM-DD format.
6. Product Statistics
   Retrieve statistics on product distribution per category and supplier.
   Method Endpoint Description
   - GET /api/statistics/products-per-category Get percentage of products per category
   - GET /api/statistics/products-per-supplier Get percentage of products per supplier

## Contribute

1. QMD <quydm.fl@gmail.com>

## Technology Stack

1. Framework: Go Gin
2. Database: PostgresQL
3. Authentication: JWT

## How to run it?

1. Clone project from github with command bellows:

```cmd
cd ${HOME}/go/src/
git clone git@github.com:quydmfl/niveau-test.git
cd niveau-test
```

2. Run docker dependencies services as: database, v.v...

```cmd
cd deploy/docker-compose
docker compose up -d
```

3. Update config database and another config

If run on local:

```cmd
vim config/local.yml
```

If run on production:

```cmd
vim config/prod.yml
```

Edit information bellows:

```
data:
  db:
    user:
      driver: postgres
      dsn: host=${HOST} user=${USER} password=${PASSWORD} dbname=${DATABASE_NAME} port=${PORT} sslmode=disable TimeZone=UTC
```

4. Run application

Install tools nunu:

```cmd
go install github.com/go-nunu/nunu@latest
```

Start application:

```cmd
nunu run
```

Ex:

```cmd
niveau-test ➤ nunu run

? Which directory do you want to run?  [Use arrows to move, type to filter]
  cmd/migration/main.go
> cmd/server/main.go
  cmd/task/main.go
```

5. Options

Migration database.

```cmd
cmd/migration/main.go
```

Start services http.

```cmd
cmd/server/main.go
```

6. Swagger documentation

```cmd
http://localhost:8000/swagger/index.html
```
