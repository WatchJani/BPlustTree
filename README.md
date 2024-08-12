# Go B+ Tree Library

[![GoDoc](https://godoc.org/github.com/yourusername/bplustree?status.svg)](https://pkg.go.dev/github.com/yourusername/bplustree)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/bplustree)](https://goreportcard.com/report/github.com/yourusername/bplustree)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A high-performance, thread-safe implementation of the B+ Tree data structure in Go. Ideal for database indexing, sorted data management, and efficient range queries.

## Features

- **Full B+ Tree Implementation**: Supports insertion, deletion, and search operations.
- **Range Queries**: Efficiently handles range queries for sorted data.
- **Customizable Node Size**: Users can define the maximum number of keys per node.
- **Iterator Support**: Provides a built-in iterator for in-order traversal.
- **Thread-Safe**: Designed for concurrent read and write operations.
- **Optimized for Performance**: Tuned for high performance in various use cases.

## Installation

To install the package, use `go get`:

```sh
go get github.com/WatchJani/BPlustTree
