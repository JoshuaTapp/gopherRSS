# gopherRSS
A RSS feed aggregator written entirely in Golang!

## What does this tool do?

This is a simple backend that does manages RSS feeds by automatically pulling them for all users and allowing frontends to consume those. 

## Features

1. Auth Middleware using API key access
2. RSS fetcher Worker using goroutines
3. Goose and sqlc to manage PostgreSQL db
