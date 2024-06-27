# fparams
___

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
go get github.com/artemk1337/fparams/cmd/main
```

### Usage

```shell

```

Parameters:
- `-disableCheckFuncParams` - disable check function params
- `-disableCheckFuncReturns` - disable check function returns


## Tests

---

Run command to start test:
```shell
go test ./pkg/analyzer/...
```

## Integrations
- golangci-lint
