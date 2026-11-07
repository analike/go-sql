# go-sql

A small Go library of helper functions for building PostgreSQL SQL queries. It helps assemble `$n` positional placeholders and the matching argument slices for `WHERE`, `INSERT`, and `UPDATE` clauses. It does not execute queries or replace a database driver.

## Requirements

- Go 1.24 or later
- A PostgreSQL driver or other query execution layer of your choice

## Install

```sh
go get go.analike.dev/sql
```

## Usage

Import the package (aliased here as `sql` to distinguish it from `database/sql`):

```go
import sql "go.analike.dev/sql"
```

### Build a WHERE clause

```go
filter := sql.NewFilter().
    WithFilter("status", "active").
    WithFilterOperator("age", 18, ">=")

query := "SELECT id, name FROM users WHERE " + filter.GetWhere()
rows, err := db.QueryContext(ctx, query, filter.GetParams()...)
```

The resulting clause is `status = $1 AND age >= $2`, and `GetParams()` returns `[]any{"active", 18}` in the same order.

Filters can also be grouped:

```go
filter := sql.NewFilter().
    WithFilter("tenant_id", tenantID).
    WithOrGroup(func(group *sql.Filter) {
        group.WithFilter("role", "admin").
        group.WithFilter("role", "owner")
    })

query := "SELECT id FROM users WHERE " + filter.GetWhere()
rows, err := db.QueryContext(ctx, query, filter.GetParams()...)
```

This builds `tenant_id = $1 AND (role = $2 OR role = $3)`. Use `WithAndGroup` or `WithGroup("AND" /* or "OR" */, build)` for other grouped conditions. Groups share the parent placeholder counter.

Other filter helpers include:

- `WithFilterReversed` and `WithFilterOperatorReversed` to put the placeholder on the left side of the comparison.
- `WithFilterOperator` to specify a comparison operator.
- `WithLike` and `WithLikeNoCase` for `LIKE` and `ILIKE` comparisons.
- `WithValueFormat` to wrap a placeholder in a caller-provided SQL format.
- `WithFilterOperatorEmbed` to add a condition containing a literal SQL value without adding a parameter.

### Build an INSERT values list

```go
insert := sql.NewInsert().WithValues("Ada", "ada@example.com")

query := "INSERT INTO users (name, email) VALUES (" + insert.GetPlaceholders() + ")"
_, err := db.ExecContext(ctx, query, insert.GetValues()...)
```

`GetPlaceholders()` returns `$1, $2`, matching the order of values supplied to `WithValue` or `WithValues`.

### Build an UPDATE SET clause

```go
update := sql.NewUpdate().
    WithValue("name", "Ada").
    WithValue("email", "ada@example.com")

query := "UPDATE users SET " + update.GetUpdates() + " WHERE id = $3"
args := append(update.GetValues(), userID)
_, err := db.ExecContext(ctx, query, args...)
```

`GetUpdates()` returns `name = $1, email = $2`. Use `WithNoSetValue` to append a parameter without adding a `SET` assignment and get its placeholder (for example, for a `WHERE` condition):

```go
update := sql.NewUpdate().WithValue("name", "Ada")
idPlaceholder := update.WithNoSetValue(userID)

query := "UPDATE users SET " + update.GetUpdates() + " WHERE id = " + idPlaceholder
_, err := db.ExecContext(ctx, query, update.GetValues()...)
```

`GetNextCount()` advances and returns the placeholder count without adding a value.

## Notes

- These helpers use PostgreSQL positional placeholders (`$1`, `$2`, …).
- Pass values through the returned argument slices; do not interpolate untrusted values into SQL text.
- Table names, column names, operators, and SQL format strings are inserted into the query text by these helpers. Only use trusted SQL fragments for those inputs; placeholders parameterize values, not identifiers or SQL syntax.
- `Filter.GetWhere()` adds the tautology `1 = $1` with parameter `"1"` when called on an empty filter.

## License

MIT. See [LICENSE](LICENSE).
