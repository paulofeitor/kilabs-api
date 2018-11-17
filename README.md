# KI Labs Challenge

Interview Calendar API

## Getting Started

This project is a Coding Challenge for KI Labs. The purpose is to create an API to manage an interview calendar with 2 roles (Interviewers and Candidates)

### Prerequisites

* [Golang](https://golang.org/) - Golang Website
* [MySQL](https://www.mysql.com/) - MySQL Website

### Installing

Run the database.sql script to create a new database and tables.
Check the config/config.go to change the Database connection info.
```go
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Driver:  "mysql",
			Host:    "localhost",
			Port:    "8889",
			User:    "root",
			Pass:    "root",
			Name:    "kilabs",
			Charset: "utf8",
		},
	}
}
```
Then just start the service:

```bash
go run main.go
```

## Structure
```
├── app
│   ├── app.go
│   ├── routes              // API routes
│   │   ├── common.go       // Common response functions
│   │   ├── candidates.go   // APIs for Candidates (CRUD)
│   │   └── interviewers.go // APIs for Interviewers (CRUD)
│   │   └── slots.go        // APIs for Slots (Matching)
│   └── model
│       └── model.go     // Structs
├── config
│   └── config.go        // Database configuration
└── main.go
```

## Built With

* [httprouter](https://github.com/julienschmidt/httprouter) - The Router Package used
Based on:
* [mingrammer REST API Example](https://github.com/mingrammer/go-todo-rest-api-example) - Example to structurize and implement routes and handlers

## Authors

* **Paulo Feitor** - [paulofeitor](https://github.com/paulofeitor)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.
