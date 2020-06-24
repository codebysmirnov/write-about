## Install
* Prepare your database (create a new database and role)
* Clone or download repo
* Go into project directory and install dependencies
```shell script
$ go mod download
```
* Install sql-migrate
```shell script
$ go get -v github.com/rubenv/sql-migrate/...
```
* Setup environment variables for your database
```shell script
export DB_HOST=<YOUR HOST>
export DB_NAME=<YOUR NAME OF DATABASE>
export DB_USER=<YOUR NAME OF ROLE>
export DB_PASSWORD=<YOUR ROLE PASSWORD>
```
* Run sql-migrate and check status of migrations
```shell script
$ sql-migrate up
$ sql-migrate status
```
* Install sqlboiler and driver for your database
```shell script
$ go get -u -t github.com/volatiletech/sqlboiler
$ go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
```
* Create sqlboiler configuration file from example (sqlboiler.toml.example)
```toml
output           = "models"
wipe             = true
no_tests         = false
no_context       = true
add_soft_deletes = true

[psql]
  dbname = "dbname"
  host   = "localhost"
  port   = 5432
  user   = "dbusername"
  pass   = "dbpassword"
  schema = "myschema"
  blacklist = ["migrations"]
```
* Run sqlboiler to generate models from your database schema
```shell script
$ sqlboiler psql --no-context --add-soft-deletes
$ go test ./<YOUR MODELS DIR>
```
