# go-lox
Bob Nystrom's Lox (from Crafting Interpreters) implemented in Go

To run all tests:
```bash
    go test ./...
```

To run the interpreter on a test file:
```bash
    go run main.go -- sample.lox
```

To build the interpreter and then run it:
```bash
    go build
    ./go-lox sample.lox
```
