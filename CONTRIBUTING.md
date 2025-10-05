# Contributing to Rejonson

First off, thank you for considering contributing to Rejonson! It's people like you that make Rejonson such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible.
* **Provide specific examples to demonstrate the steps**.
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**
* **Include details about your configuration and environment**.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
* **Provide specific examples to demonstrate the steps**.
* **Describe the current behavior** and **explain which behavior you expected to see instead** and why.
* **Explain why this enhancement would be useful** to most Rejonson users.

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

### Setting Up Your Development Environment

1. Fork the repo
2. Clone your fork
3. Install Go 1.21 or later
4. Install Docker (for running Redis with ReJSON module)
5. Run tests: `bash test.sh`

### Code Generation

This project uses code generation to support multiple versions of go-redis. If you modify the generator:

1. Navigate to the `generator` directory
2. Run `go run .` to regenerate code
3. Test all versions

### Testing

Before submitting a PR, ensure:

1. All tests pass: `bash test.sh`
2. Code is formatted: `go fmt ./...`
3. No linting issues: `go vet ./...`
4. All versions build successfully

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

## Project Structure

```
.
├── generator/          # Code generator for multiple go-redis versions
├── v6.generated.go     # Generated code for go-redis v6
├── v7/                 # go-redis v7 support
├── v8/                 # go-redis v8 support
└── v9/                 # go-redis v9 support
```

## Questions?

Feel free to open an issue with your question, and we'll do our best to help!

Thank you for contributing! 🎉
