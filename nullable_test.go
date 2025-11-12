package nullable_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zikwall/nullable"
)

func TestNew(t *testing.T) {
	n := nullable.New(42)
	require.False(t, n.IsNull())
	require.True(t, n.NotNull())
	require.Equal(t, 42, n.Value())
}

func TestNullConstructor(t *testing.T) {
	n := nullable.Null[string]()
	require.True(t, n.IsNull())
	require.False(t, n.NotNull())
	require.Empty(t, n.Value())
	require.Nil(t, n.Ref())

	b, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, "null", string(b))
}

func TestFromRef(t *testing.T) {
	val := "hello"
	n := nullable.FromRef(&val)
	require.True(t, n.NotNull())
	require.Equal(t, "hello", n.Value())

	var nilPtr *string
	n2 := nullable.FromRef(nilPtr)
	require.True(t, n2.IsNull())
	require.Empty(t, n2.Value())
}

func TestConvertRef(t *testing.T) {
	type user struct{ Name string }
	u := &user{"Bob"}

	n := nullable.ConvertRef(u, func(u *user) string { return u.Name })
	require.True(t, n.NotNull())
	require.Equal(t, "Bob", n.Value())

	var nilUser *user
	n2 := nullable.ConvertRef(nilUser, func(u *user) string { return u.Name })
	require.True(t, n2.IsNull())
}

func TestValue_IsNull_NotNull(t *testing.T) {
	n := nullable.Nullable[int]{}
	require.True(t, n.IsNull())
	require.False(t, n.NotNull())
	require.Equal(t, 0, n.Value())

	n2 := nullable.New(99)
	require.False(t, n2.IsNull())
	require.Equal(t, 99, n2.Value())
}

func TestRef(t *testing.T) {
	n := nullable.New("data")
	ref := n.Ref()
	require.NotNil(t, ref)
	require.Equal(t, "data", *ref)

	n2 := nullable.Nullable[string]{}
	require.Nil(t, n2.Ref())
}

func TestSetAndUnset(t *testing.T) {
	n := nullable.Null[int]()
	require.True(t, n.IsNull())

	n.Set(10)
	require.True(t, n.NotNull())
	require.Equal(t, 10, n.Value())

	n.Unset()
	require.True(t, n.IsNull())
	require.Equal(t, 0, n.Value())
}

func TestStringMethod(t *testing.T) {
	n := nullable.New("hello")
	require.Equal(t, "hello", n.String())

	n2 := nullable.Null[string]()
	require.Equal(t, "null", n2.String())
}

func TestEqualMethod(t *testing.T) {
	n1 := nullable.New(42)
	n2 := nullable.New(42)
	n3 := nullable.New(100)
	nNull := nullable.Null[int]()

	require.True(t, n1.Equal(n2), "equal values should match")
	require.False(t, n1.Equal(n3), "different values should not match")
	require.False(t, n1.Equal(nNull), "non-null vs null should differ")
	require.True(t, nNull.Equal(nullable.Null[int]()), "both nulls should be equal")
}

func TestEqualWithStruct(t *testing.T) {
	type Point struct{ X, Y int }
	a := nullable.New(Point{1, 2})
	b := nullable.New(Point{1, 2})
	c := nullable.New(Point{2, 3})

	require.True(t, a.Equal(b))
	require.False(t, a.Equal(c))
}

func TestSetThenMarshal(t *testing.T) {
	n := nullable.Null[string]()
	n.Set("go")
	require.True(t, n.NotNull())

	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, `"go"`, string(data))
}

func TestMarshalJSON_NotNull(t *testing.T) {
	n := nullable.New("abc")
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, `"abc"`, string(data))
}

func TestMarshalJSON_Null(t *testing.T) {
	n := nullable.Nullable[int]{}
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, "null", string(data))
}

func TestUnmarshalJSON_Value(t *testing.T) {
	var n nullable.Nullable[int]
	err := json.Unmarshal([]byte("123"), &n)
	require.NoError(t, err)
	require.True(t, n.NotNull())
	require.Equal(t, 123, n.Value())
}

func TestUnmarshalJSON_Null(t *testing.T) {
	var n nullable.Nullable[int]
	err := json.Unmarshal([]byte("null"), &n)
	require.NoError(t, err)
	require.True(t, n.IsNull())
	require.Equal(t, 0, n.Value())
}

func TestUnmarshalJSON_Empty(t *testing.T) {
	var n nullable.Nullable[string]
	err := json.Unmarshal([]byte("null"), &n)
	require.NoError(t, err)
	require.True(t, n.IsNull())
}

func TestUnmarshalJSON_Empty2(t *testing.T) {
	var n nullable.Nullable[string]

	err := n.UnmarshalJSON([]byte(""))
	assert.NoError(t, err)
	assert.True(t, n.IsNull())
}

func TestIsZero(t *testing.T) {
	n := nullable.Nullable[int]{}
	require.True(t, n.IsZero())

	n2 := nullable.New(5)
	require.False(t, n2.IsZero())
}
