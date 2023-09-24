# Product API

## Implements:
-- Full testing
-- Modeling
-- Gorilla/Mux

## Clone the repo
```bash
$ git clone github.com/DLzer/go-product-api
```

## Testing Database
We'll set up a test instance of PostgreSQL using docker for convenience using the following command
```bash
$ docker run --name postgres-docker -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```
This will pull the Postgres image from docker, set the default password, the name for the container
as well as setting the port to `5432`.

## Set your default PostreSQL variables
If you're lucky enough to be working in a *nix system, run these to set the default environment variables
for connecting to your test database.
```bash
export APP_DB_USERNAME=postgres
export APP_DB_PASSWORD=postgres
export APP_DB_NAME=postgres
```
Otherwise for windows, the easiest method is setting them via the environment variable manager. 
*Note* 1. Set the variables for the USER scope, not for system. 2. You'll have to log out, and back in
for the changes to take effect.

## Running the tests
Run the test suite with the following commands
```bash
$ go test -v
```
The expected output should be:
```bash
=== RUN   TestEmptyTable
--- PASS: TestEmptyTable (0.01s)
=== RUN   TestGetNonExistentProduct
--- PASS: TestGetNonExistentProduct (0.00s)
=== RUN   TestCreateProduct
--- PASS: TestCreateProduct (0.01s)
=== RUN   TestGetProduct
--- PASS: TestGetProduct (0.01s)
=== RUN   TestUpdateProduct
--- PASS: TestUpdateProduct (0.01s)
=== RUN   TestDeleteProduct
--- PASS: TestDeleteProduct (0.01s)
PASS
ok      github.com/DLzer/go-product-api       0.071s
```

## Deployment
This application comes with dockerfiles prepped for deployment of the application, postres, and PGAdmin.
With docker installed and running - run the command:
```bash
$ docker-compose up --build
```
This will fire up a new app with all three containers. After it's up an running just visit `http://localhost:8080` 
and you'll be greeted accordingly!


## Version
1.0.0
