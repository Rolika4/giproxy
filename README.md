[![codecov](https://codecov.io/github/Rolika4/giproxy/branch/main/graph/badge.svg?token=4G13TLABQY)](https://codecov.io/github/Rolika4/giproxy)

# Giproxy

Giproxy is a lightweight proxy server designed to facilitate interactions with Git providers (like GitLab, GitHub, and Bitbucket). The tool simplifies and unifies API requests, making integration with various Git services effortless.

---

## Features

- Unified interface for multiple Git providers
- Lightweight and fast
- Extensible architecture for new providers
- Easy to use with minimal configuration

---

## Development

To set up the development environment and start working on Giproxy, follow these steps:

### Prerequisites

- [Golang](https://go.dev/doc/install) (1.20+ recommended)
- [Node.js](https://nodejs.org/en/) and [npm](https://www.npmjs.com/)
- [Make](https://www.gnu.org/software/make/)
- Git

### Install Dependencies

Run the following command to install the necessary dependencies:

```bash
make install
```

This will ensure all Go modules are downloaded and ready.

### Developer Mode

To start the project in developer mode with hot-reloading:

```bash
make watch
```

This will automatically restart the application whenever changes are detected in the source code.

### Standard Run

If you want to manually start the application without hot-reloading:

```bash
make run
```

---

## Build

To build the project and generate an executable binary:

```bash
make build
```

The built binary will be available in the `bin/` directory.

---

## API Documentation

API documentation for Giproxy is currently under development. Stay tuned for updates.

---

## Contribution

Contributions are welcome! Please follow the guidelines below:

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Open a pull request