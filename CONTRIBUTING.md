# Contributing

To work on the project, install the following versions of the tools.

```sh
$ cat .tool-versions
go            1.21.5
golangci-lint 1.55.2
goreleaser    1.24.0
# make        4.4.1
```

Once you have installed the tools, download the dependencies.

```sh
go mod download
```

Now you are ready to start working.

```sh
$ make
all       Run all recipes.
build     Build a binary.
help      Show help information.
lint      Lint the source code.
test      Run tests.
```

Thank you for taking the time to read this. Good luck.
