# pgx-google-uuid

[`github.com/google/uuid`](https://github.com/google/uuid) type support
for [`github.com/jackc/pgx`](https://github.com/jackc/pgx) PostgreSQL driver

Major package version corresponds to the major pgx version, e.g.:

- `github.com/vgarvardt/pgx-google-uuid/v4` -> `github.com/jackc/pgx/v4`
- `github.com/vgarvardt/pgx-google-uuid/v5` -> `github.com/jackc/pgx/v5`

## Usage example

```go
package main

import (
  "context"
  "os"

  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func main() {
  pgxConfig, err := pgxpool.ParseConfig(os.Getenv("PG_URI"))
  if err != nil {
    panic(err)
  }

  pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    pgxUUID.Register(conn.TypeMap())
    return nil
  }

  pgxConnPool, err := pgxpool.ConnectConfig(context.TODO(), pgxConfig)
  if err != nil {
    panic(err)
  }

  // use pgxConnPool
  ...
}
```
