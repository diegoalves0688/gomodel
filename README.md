# digital-go-sample

This is a sample project to create Go based products on Digital.

## Main dependencies

### Packages

The main packages used in the project:

[Bun DB Client](https://bun.uptrace.dev/): It is *sql.DB compatible, database-agnostic and supports fixtures. Bun is a rewrite of go-pg (maintenance mode) and slightly less efficient but works with different databases.

[Echo HTTP Router](https://echo.labstack.com/): It presents a good benchmark compared to Gin https://github.com/labstack/echo#benchmarks and allows the approach to centralized manipulation of errors.

[Echo Swagger](https://github.com/swaggo/echo-swagger): It is a echo middleware to automatically generate RESTful API documentation with Swagger 2.0.

[Viper Config](https://github.com/spf13/viper): It can manipulate all types of configuration needs and formats. Here it was good to work with a profile config strategy, using a .yml for each environment and supporting merging of config files.

[Tescontainers](https://github.com/testcontainers/testcontainers-go): It is simple to handle container-based dependencies for automated integration tests.

### Development tools

The main tools used in the project:

[Migrate](https://github.com/golang-migrate/migrate): It is a well-known library in the Go community, and for our purpose of running migrations via CLI it works fine.

[Swag](https://github.com/swaggo/swag): it generates RESTful API documentation with Swagger 2.0. 

[Mockery](https://github.com/vektra/mockery): it allows to generate mocks for interfaces and removes the boilerplate code required to use mocks.

[Golangci-lint](https://github.com/golangci/golangci-lint): it runs linters in parallel, uses caching, supports yaml config, has integrations with all major IDE and has dozens of linters included.

## Installing development dependencies (tools)

The project uses Task because it is easier and simpler than Make and allows to create a good semantic for tasks, for example, we can have tasks for context: test:unit, test:integration. Before installing dependencies you need to install Task (if you want another way, you can find [here](https://taskfile.dev/#/installation)):

```go install github.com/go-task/task/v3/cmd/task@latest```

Now you can install the development dependencies:

```task install:dev-deps```

## Running local

After installing the development dependencies, you can run the project in your machine:

Starting the services such as database:

```docker-compose up```

Running the application:

```task gosample:run```

## Available tasks

**TaskInstall** - holds installation tasks:

> Installing dependencies: ```task install:dev-deps```

**TaskLint** - holds lint tasks:

> Running all lints: ```task lint:run```

**TaskTests** - holds test tasks:

> Running unit tests: ```task test:unit```

> Running unit tests with coverage: ```task test:unit-coverage```

> Running integration tests: ```task test:integration```

> Generating mocks: ```task test:gen-mocks```

## About fixtures

They are a good resource for integration tests because allow inserting data before the test execution. Here we are using the fixtures feature provided for Bun.

```
 func (s *MessageHttpIntegrationTest) Test_GetAllMessages() {
    var models []interface{}
	models = append(models, (*domain.Message)(nil))

	s.Run(
		th.Fixture(s.Ctx, models, "message.yml"),
		func(db *bun.DB) {
			// your test...
		},
	)
}
```

The file message.yml is in the folder ./testdata, and it has the data needed to run the integration test.

```
- model: Message
  rows:
    - _id: bf40e5e7-e3a0-4a11-b1e6-3b2b0455a9ae
      receiver: Paulo
      sender: Maria
      content: 'sample'
```