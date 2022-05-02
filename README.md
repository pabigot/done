# done

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Doc](https://pkg.go.dev/badge/github.com/pabigot/done.svg)](https://pkg.go.dev/github.com/pabigot/done)
[![Go Report Card](https://goreportcard.com/badge/github.com/pabigot/done)](https://goreportcard.com/report/github.com/pabigot/done)
[![Build Status](https://github.com/pabigot/done/actions/workflows/core.yml/badge.svg)](https://github.com/pabigot/done/actions/workflows/core.yml)
[![Coverage Status](https://coveralls.io/repos/github/pabigot/done/badge.svg)](https://coveralls.io/github/pabigot/done)

Package done provides an concurrency-safe implementation of the Done() and
Err() methods as used in context.Context to communication completion status
to an application.

Why?  Because I have a lot of active objects and want a uniform way of
confirming that they've completed.
