# First Golang server

## How to run

- Create a new database (SQL Server):

```
CREATE DATABASE [simple_server];
```

- Then, create new table in this database:

```
CREATE TABLE dbo.Users
(
    UserId INT IDENTITY(1,1) PRIMARY KEY,
    Username VARCHAR(50) NOT NULL,
    Password VARCHAR(120) NOT NULL,
    Email VARCHAR(50) NOT NULL
);
```

- Inside **main.go**, replace your own credentials:

```
var server = "localhost"
var port = 1433
var user = "sa"
var password = "your_password"
var database = "simple_server"
```

- Finally, run as:

```
go run main.go
```