# fparams

### Overview

fparams is a linter tool for Go that checks the formatting of function parameters and return values. 
The linter ensures that function parameters and return values are either all on one line or each on a new line. 
This helps maintain consistent and readable code formatting in Go projects.


### Features

- Parameter and return value checking: Verifies if function parameters and return values are formatted correctly.
- Configurable checks: Allows disabling checks for function parameters or return values through flags.
- Automated suggestions: Provides suggested fixes to format parameters and return values on separate lines if needed.

### Instalation

```shell
go install github.com/artemk1337/fparams/cmd/fparams@latest
```

### Usage

```shell
fparams ./...
```

### Configuration

You can configure fparams using command-line flags to enable or disable specific checks:
- `-disableCheckFuncParams` - disable check function params
- `-disableCheckFuncReturns` - disable check function returns

### Example

Given a function declaration like this:
```go
func example(a int, b int, 
    c, d string) (int, error) {
    return 0, nil
}
```

fparams will suggest changing it to:
```go
func example(
    a int,
    b int,
    c string,
    d string,
) (
    int,
    error,
) {
    return 0, nil
}
```

More valid and invalid examples see in [testdata](pkg%2Fanalyzer%2Ftestdata) directory.

### Tests

Run command to start test:
```shell
go test ./pkg/analyzer/...
```

[//]: # (### Integrations)
[//]: # (- golangci-lint)

### Contribution

Contributions are welcome! 
Please fork the repository and submit a pull request.

### License

This project is licensed under the MIT License. 
See the [LICENSE](LICENSE) file for details.

By using fparams, you can ensure that your Go codebase remains clean and consistently formatted, 
making it easier to read and maintain.