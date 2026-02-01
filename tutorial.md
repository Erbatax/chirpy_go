## Migrate database

go to the sql schema directory

```bash
cd sql/schema
```

### Up

```bash
goose postgres "postgres://user:password@localhost:5432/chirpy" up
```

### Down

```bash
goose postgres "postgres://user:password@localhost:5432/chirpy" down
```

## Create code from sql queries

```bash
sqlc generate
```

## Run server

```bash
go run .
```
