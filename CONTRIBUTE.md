# Contributing to dot-agent

Thank you for your interest in contributing to `dot-agent`! We welcome contributions from everyone.

This document provides guidelines for contributing to the project. These are just guidelines, not rules. Use your best judgment and feel free to propose changes to this document in a pull request.

## How Can I Contribute?

### Reporting Bugs

Bugs are tracked as [GitHub issues](https://github.com/cthulhu/dot-agent/issues). When creating a bug report, please include as many details as possible:

*   **Use a clear and descriptive title** for the issue to identify the problem.
*   **Describe the exact steps which reproduce the problem** in as many details as possible.
*   **Provide specific examples to demonstrate the steps.**
*   **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
*   **Explain which behavior you expected to see instead and why.**
*   **Include screenshots and animated GIFs** which help you demonstrate the steps or the bug.
*   **Include your environment details** (OS, Go version, `dot-agent` version).

### Suggesting Enhancements

Enhancement suggestions are also tracked as [GitHub issues](https://github.com/cthulhu/dot-agent/issues).

*   **Use a clear and descriptive title** for the issue to identify the suggestion.
*   **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
*   **Explain why this enhancement would be useful** to most `dot-agent` users.
*   **List some other tools or applications where this feature exists**, if applicable.

### Contributing Code

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** to your local machine.
3.  **Create a new branch** for your work (e.g., `git checkout -b feature/my-cool-feature` or `git checkout -b fix/issue-123`).
4.  **Make your changes.** Ensure you follow the project's coding standards and add tests for new functionality or bug fixes.
5.  **Run the tests** to make sure everything is working as expected: `go test ./...`.
6.  **Commit your changes.** Use clear and descriptive commit messages.
7.  **Push your branch** to your fork on GitHub.
8.  **Submit a pull request** to the `main` branch of the original repository.

## Development Setup

### Prerequisites

*   [Go 1.22+](https://golang.org/dl/)
*   [Git](https://git-scm.com/)

### Getting Started

1.  Clone the repository:
    ```bash
    git clone https://github.com/cthulhu/dot-agent.git
    cd dot-agent
    ```

2.  Download dependencies:
    ```bash
    go mod download
    ```

3.  Build the project:
    ```bash
    go build -o dot-agent ./cmd/dot-agent
    ```

4.  Run tests:
    ```bash
    go test ./...
    ```

## Coding Standards

*   **Formatting**: Use standard Go formatting. Run `go fmt ./...` before committing.
*   **Testing**: We aim for high test coverage. Please add unit tests for any new features or bug fixes.
*   **Documentation**: Update the `README.md` if you add new commands or change existing behavior.
*   **Simplicity**: Keep code simple and easy to understand. Follow standard Go idioms.

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
