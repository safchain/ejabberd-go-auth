# ejabberd-go-auth
Ejabberd External Auth in GO

## Build

make

## Configuration

```sh
Driver = postgres
Host = 172.17.0.1
Port = 5432
User = postgres
Pass = postgres
Dbname = postgres
Dbargs = sslmode=disable
Table = users
UserField = user
PassField = passwd
```

## Run

ejabberd-go-auth -conf ejabberd-go-auth.ini

## Operations supported

* isuser
* auth

## Backends supported

* PostgreSQL
* MySQL
