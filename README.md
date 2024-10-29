[![build](https://github.com/zikwall/nullable/workflows/golangci_lint/badge.svg)](https://github.com/zikwall/nullable/actions)

<div align="center">
  <h1>Nullable</h1>
</div>

## Install

- `$ go get -u github.com/zikwall/nullable`

## Usage

```go
type SomeStruct struct {
	DeviceID nullable.Nullable[string]    `json:"device_id"`
	GUID     nullable.Nullable[uuid.UUID] `json:"-"`

	DeviceType  nullable.Nullable[string] `json:"device_type"`
	DeviceModel nullable.Nullable[string] `json:"device_model"`

	OSName    nullable.Nullable[string] `json:"os_name"`
	OSVersion nullable.Nullable[string] `json:"os_version"`

	UserAgent      nullable.Nullable[string] `json:"user_agent"`
	Browser        nullable.Nullable[string] `json:"browser"`
	BrowserVersion nullable.Nullable[string] `json:"browser_version"`

	UserID        nullable.Nullable[int64] `json:"user_id"`

	LocalTimestamp nullable.Nullable[time.Time] `json:"local_timestamp"`
}

some := SomeStruct{}
some.GUID = nullable.New[uuid.UUID](uuid.New())

fmt.Println(some.GUID.NotNull()) // true
fmt.Println(some.GUID.IsNull()) // false

v := 42.42
var ref *int64 = &v
var ref2 *int64

nullableRef := nullable.FromRef(ref)
nullableRef2 := nullable.FromRef(nil)

fmt.Println(nullableRef.NotNull()) // true
fmt.Println(nullableRef2.IsNull()) // true
```