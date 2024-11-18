# Go Core Library

[![GoDoc](https://pkg.go.dev/badge/github.com/tittuvarghese/go-core-library)](https://pkg.go.dev/github.com/tittuvarghese/go-core-library)
[![Build Status](https://travis-ci.org/tittuvarghese/go-core-library.svg?branch=main)](https://travis-ci.org/tittuvarghese/go-core-library)

The Go Core Library provides a set of foundational components designed for use in any Go service or application. This library includes modules for configuration management, cryptography (hashing), JWT handling, logging, storage (both in-memory and database), time utilities, and validation. It aims to simplify the common setup for Go services and make your projects more modular and maintainable.

## Features

- **Config Management**: Simplified access and management of configuration settings, supporting multiple backends (e.g., environment variables, JSON, YAML, etc.).
- **Crypto**: Provides utilities for cryptographic operations like hashing (e.g., SHA256, bcrypt).
- **JWT**: Create, sign, and validate JWT tokens for authentication and authorization.
- **Logger**: Integrated logging capabilities with support for different log levels and formats.
- **Storage**: Abstraction over database and in-memory storage for easy integration with various data backends.
- **Time**: A wrapper around time-related functions to facilitate unit testing and ensure consistent time handling.
- **Validator**: Input validation utilities to ensure data integrity and prevent errors.

## Installation

To install the Go Core Library, run the following command:

```bash
go get github.com/tittuvarghese/ss-go-core

