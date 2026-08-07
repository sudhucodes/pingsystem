# Contributing to PingSystem

Thank you for your interest in contributing to **PingSystem**! We welcome contributions from the community to help improve this Windows background monitoring agent.

---

## 📜 Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

---

## 🛠️ Development Environment Setup

### Prerequisites

-   **Go**: 1.26 or newer ([golang.org](https://golang.org/))
-   **Git**: [git-scm.com](https://git-scm.com/)

### Getting Started

1. **Fork the Repository** on GitHub.
2. **Clone your fork**:
    ```bash
    git clone https://github.com/sudhucodes/pingsystem.git
    cd pingsystem
    ```
3. **Run Tests Locally**:
    ```bash
    go test ./...
    ```
4. **Verify Cross-Compilation**:
    ```bash
    ./build.sh
    # Or manually:
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o PingSystem.exe .
    ```

---

## 📦 Release Process

PingSystem releases are automatically managed based on the `version` field in `package.json`.

When a maintainer bumps the `"version"` in `package.json` (e.g. `"0.0.2"` -> `"0.0.3"`) on the `main` branch:
1. GitHub Actions automatically creates and pushes the release tag `v0.0.3`.
2. Cross-compiles the native Windows binary `PingSystem.exe`.
3. Creates a ZIP package containing the binary and license.
4. Publishes a GitHub Release with `PingSystem.exe` directly attached.

---

## 🧹 Code Quality & Standards

Before committing your changes, ensure:

-   Code is properly formatted:
    ```bash
    go fmt ./...
    ```
-   Go vet checks pass cleanly:
    ```bash
    go vet ./...
    ```
-   All unit tests pass:
    ```bash
    go test ./...
    ```

---

## 🚀 Submitting a Pull Request

1. Create a feature branch off `main`:
    ```bash
    git checkout -b feature/my-cool-feature
    ```
2. Commit your changes with clear, descriptive commit messages.
3. Push to your fork and open a Pull Request against `main`.
4. Ensure all GitHub CI checks pass!

---

## 📄 License

By contributing to PingSystem, you agree that your contributions will be licensed under the project's [Apache License 2.0](LICENSE).
