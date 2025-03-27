# Feed Me

Feed me is a humble platform to manage IP, URL or domain feeds.

It is designed to create feeds and to ingest

## Installation

First of all, you need to create a database.

```bash
sudo -u postgres psql
```

```sql
CREATE USER __user__ WITH ENCRYPTED PASSWORD 'changeme';
CREATE DATABASE feedme;
GRANT ALL PRIVILEGES ON feedme TO __user__;
```

## Development

As an Go project, you first need to install Golang. It can be found in most Linux distributions.

It algo uses Sqlc to generate Go code from the SQL Statements.

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
