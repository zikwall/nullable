[![build](https://github.com/zikwall/nullable/workflows/golangci_lint/badge.svg)](https://github.com/zikwall/nullable/actions)

<div align="center">
  <h1>Nullable</h1>
</div>

nullable provides a *type-safe and JSON-friendly way* to represent optional (nullable) values in Go.
It eliminates unsafe pointer juggling and keeps your models clean, especially when dealing with APIs or databases.

## ✨ Features

- ✅ Type-safe — no interface{} or reflection
- 🧠 Generic — works with any type (int, string, uuid.UUID, structs…)
- 🔄 JSON-aware — encodes and decodes null seamlessly
- 💾 DB-friendly — ideal for Postgres jsonb or SQL NULL values
- 🧰 Minimal API — New, Null, FromRef, Set, Unset, IsNull, Value, Ref

## 📘 Comparison

| Approach                          | Drawbacks                                  | nullable advantage                      |
| --------------------------------- | ------------------------------------------ | --------------------------------------- |
| `*T` pointers                     | risk of nil dereference, not JSON-friendly | safe, JSON `null` handled automatically |
| `sql.NullString`, `sql.NullInt64` | separate type per kind                     | one generic type for all                |
| `interface{}`                     | loses type safety                          | strongly typed generics                 |


## Install

- `$ go get -u github.com/zikwall/nullable`

## 🚀 Usage

Working with optional values in Go often ends up with *T, sql.NullString, or awkward omitempty tricks.
nullable fixes that — it gives you a clean, generic, and JSON-safe way to express “this field might be missing”.

```go
package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"yourmodule/nullable"
)

type DeviceInfo struct {
	DeviceID        nullable.Nullable[string]    `json:"device_id"`
	GUID            nullable.Nullable[uuid.UUID] `json:"guid"`
	DeviceType      nullable.Nullable[string]    `json:"device_type"`
	OSName          nullable.Nullable[string]    `json:"os_name"`
	UserID          nullable.Nullable[int64]     `json:"user_id"`
	LastSeenAt      nullable.Nullable[time.Time] `json:"last_seen_at"`
}

func main() {
	var info DeviceInfo

	// Set some values
	info.GUID = nullable.New(uuid.New())
	info.DeviceType = nullable.New("mobile")

	fmt.Println(info.GUID.NotNull()) // true
	fmt.Println(info.DeviceType.Value()) // "mobile"

	// Create from pointer reference
	val := int64(42)
	num := nullable.FromRef(&val)
	fmt.Println(num.Value()) // 42

	// Unset values explicitly
	num.Unset()
	fmt.Println(num.IsNull()) // true

	// Convert to/from JSON safely
	data, _ := json.Marshal(info)
	fmt.Println(string(data))
}
```

Output:

```json
{
  "device_id": null,
  "guid": "c8f8ad4a-6e21-4e6b-a8b1-fb99d5ce1d7b",
  "device_type": "mobile",
  "os_name": null,
  "user_id": null,
  "last_seen_at": null
}
```