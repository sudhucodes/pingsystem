# Contributing to PingSystem

Thank you for your interest in contributing to **PingSystem**! We welcome contributions from the community to help improve this Windows background monitoring agent.

---

## 📜 Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

---

## 🛠️ Development Environment Setup

### Prerequisites

- **Go**: 1.26 or newer ([golang.org](https://golang.org/))
- **Node.js**: 18+ and `npm` (required for `@changesets/cli` release versioning)
- **Git**: [git-scm.com](https://git-scm.com/)

### Getting Started

1. **Fork the Repository** on GitHub.
2. **Clone your fork**:
   ```bash
   git clone https://github.com/sudhucodes/pingsystem.git
   cd pingsystem
   ```
3. **Install Node Dependencies** (for changeset tooling):
   ```bash
   npm install
   ```
4. **Run Tests Locally**:
   ```bash
   go test ./...
   ```
5. **Verify Cross-Compilation**:
   ```bash
   ./build.sh
   # Or manually:
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o PingSystem.exe .
   ```

---

## 📦 Versioning & Changesets Workflow

PingSystem uses [Changesets](https://github.com/changesets/changesets) to automate semantic versioning and changelog generation.

### Adding a Changeset to Your PR

Whenever you submit a pull request that introduces a fix, feature, or breaking change:

1. Run the changeset CLI:
   ```bash
   npx changeset
   ```
2. Select the type of bump (`patch`, `minor`, `major`):
   - **patch**: Bug fixes, documentation updates, non-breaking internal cleanup.
   - **minor**: New features or expanded event hooks (backwards-compatible).
   - **major**: Breaking changes to CLI flags or configuration schema.
3. Write a concise summary of the change.
4. Commit the generated markdown file in `.changeset/` along with your PR.

---

## 🧹 Code Quality & Standards

Before committing your changes, ensure:

- Code is properly formatted:
  ```bash
  go fmt ./...
  ```
- Go vet checks pass cleanly:
  ```bash
  go vet ./...
  ```
- All unit tests pass:
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
3. Include a changeset (`npx changeset`).
4. Push to your fork and open a Pull Request against `main`.
5. Ensure all GitHub CI checks pass!

---

## 📄 License

By contributing to PingSystem, you agree that your contributions will be licensed under the project's [Apache License 2.0](LICENSE).
