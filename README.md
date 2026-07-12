# go-env-parser

Explicit, batch-validated parsing of environment variables into typed Go values.

```bash
go get github.com/cblauvelt/go-env-parser
```

## Why

Most env-parsing libraries either bind struct tags (implicit, hard to grep) or return an error per call (fail on the first problem). This library does neither:

- Every variable is read with an **explicit function call** — the mapping between env var and field is visible at the call site.
- A `Parser` **accumulates errors** as it goes. One `Err()` call at the end reports every missing or malformed variable in a single message.
- Values are read through a **`Lookuper` interface**, so tests inject a `MapLookuper` directly instead of reaching for `os.Setenv` / `t.Setenv`.

## Quick start

```go
import env "github.com/cblauvelt/go-env-parser"

p := env.New()
cfg := Config{
    DBHost:  p.RequiredString("DB_HOST"),
    DBPort:  p.Int("DB_PORT", 5432),
    Debug:   p.Bool("DEBUG", false),
    Timeout: p.Duration("TIMEOUT", 30*time.Second),
}
if err := p.Err(); err != nil {
    log.Fatal(err) // all problems reported in one message
}
```

If `DB_HOST` is unset and `DB_PORT` contains `"abc"` (a typo), the error from `p.Err()` names both problems:

```bash
env "DB_HOST" is required but not set
env "DB_PORT": invalid int value "abc": ...
```

## Required vs defaulted accessors

Every type has two families:

| Family                                             | Behavior on absent/malformed                                                                                    |
| -------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `String`, `Int`, `Bool`, …                         | Returns `def`. Never records an error. Malformed present values fire the [bad-default hook](#bad-default-hook). |
| `RequiredString`, `RequiredInt`, `RequiredBool`, … | Records a missing or parse error. Returns the zero value of `T`.                                                |

## Supported types

### Scalars

| Method                              | Type                                                                    |
| ----------------------------------- | ----------------------------------------------------------------------- |
| `String` / `RequiredString`         | `string`                                                                |
| `Int` / `RequiredInt`               | `int`                                                                   |
| `Int32` / `RequiredInt32`           | `int32`                                                                 |
| `Int64` / `RequiredInt64`           | `int64`                                                                 |
| `Uint` / `RequiredUint`             | `uint`                                                                  |
| `Uint64` / `RequiredUint64`         | `uint64`                                                                |
| `Float64` / `RequiredFloat64`       | `float64`                                                               |
| `Bool` / `RequiredBool`             | `bool` (accepts `1/0`, `true/false`, `t/f`, `TRUE/FALSE`, `True/False`) |
| `Duration` / `RequiredDuration`     | `time.Duration` (e.g. `"300ms"`, `"1.5h"`)                              |
| `Time` / `RequiredTime`             | `time.Time` (RFC 3339)                                                  |
| `TimeLayout` / `RequiredTimeLayout` | `time.Time` (custom layout)                                             |

### Collections

```go
// []string — comma-separated, each element trimmed, empty elements dropped
origins := p.CSV("CORS_ORIGINS", []string{"*"})

// alternate separator
tags := p.CSVSep("TAGS", ":", nil)

// required: error if the key is absent (an empty value is not an error)
hosts := p.RequiredCSV("ALLOWED_HOSTS")

// []int
ports := p.IntSlice("PORTS", []int{8080})
ports := p.RequiredIntSlice("PORTS")
```

### Validators

**`OneOf`** — enum membership:

```go
level := p.OneOf("LOG_LEVEL", "info", "debug", "info", "warn", "error")
level := p.RequiredOneOf("LOG_LEVEL", "debug", "info", "warn", "error")
```

**`Validated[T]`** — generic escape hatch for any parse+validate logic:

```go
parseURL := func(s string) (*url.URL, error) {
    u, err := url.Parse(s)
    if err != nil || u.Scheme == "" {
        return nil, fmt.Errorf("must be an absolute URL")
    }
    return u, nil
}

u := env.RequiredValidated(p, "WEBHOOK_URL", parseURL)
u := env.Validated(p, "WEBHOOK_URL", defaultURL, parseURL)
```

## Cross-variable constraints

Some rules span more than one variable. These three helpers check pure presence across a set of keys, append at most one error each into the same `Err()` result as every accessor, and respect [`WithEmptyAsUnset`](#withemptyasunset) (a set-but-empty variable counts as unset). They never panic on zero or one key.

**`AllOrNone`** — the keys must be all set or all unset. Use it for values that only make sense together, such as a certificate and its private key:

```go
p.AllOrNone("TLS_CERT", "TLS_KEY")
```

**`MutuallyExclusive`** — at most one of the keys may be set. Use it for competing sources of the same thing:

```go
p.MutuallyExclusive("DB_DSN", "LOCAL_STORE")
```

**`RequiredWith`** — when the first key is set, every dependent must also be set. Use it for variables that a feature key turns from optional into required (one error per missing dependent):

```go
p.RequiredWith("SESSION_KEY", "API_URL", "API_KEY")
```

## Options

### `WithEmptyAsUnset`

Treat `KEY=""` the same as an absent key. Off by default, so explicit-empty values are preserved.

```go
p := env.New(env.WithEmptyAsUnset())
```

### Bad-Default Hook

Called when a **defaulted** accessor finds a present-but-malformed value. Useful for logging typos that would otherwise be silently replaced by the default.

```go
p := env.New(env.WithBadDefaultHook(func(key, raw string, err error) {
    slog.Warn("env var malformed, using default", "key", key, "raw", raw, "err", err)
}))
```

## Testing

Inject a `MapLookuper` so tests never touch the real environment:

```go
func TestLoadConfig(t *testing.T) {
    p := env.From(env.MapLookuper{
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DEBUG":   "true",
    })
    cfg, err := LoadConfig(p)
    // ...
}
```

A function that accepts a `*env.Parser` can be tested this way without `os.Setenv`, `t.Setenv`, or any global state.

## Porting from explicit `os.Getenv` calls

Before:

```go
var missing []string
dbURL := func(key string) string {
    v := os.Getenv(key)
    if v == "" { missing = append(missing, key) }
    return v
}("DB_URL")
port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
if port == 0 { port = 5432 }
if len(missing) > 0 {
    return fmt.Errorf("missing: %s", strings.Join(missing, ", "))
}
```

After:

```go
p := env.New()
dbURL := p.RequiredString("DB_URL")
port  := p.Int("DB_PORT", 5432)
if err := p.Err(); err != nil {
    return err
}
```
