package nullable

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		want := Nullable[int]{
			value:   7,
			notNull: true,
		}

		assert.Equal(t, want, New(7))
	})

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()

		want := Nullable[[]int]{
			value:   nil,
			notNull: true,
		}

		var slice []int

		assert.Equal(t, want, New(slice))
		assert.Equal(t, want, New[[]int](nil))
	})

	t.Run("any", func(t *testing.T) {
		t.Parallel()

		want := Nullable[any]{
			value:   nil,
			notNull: true,
		}

		var v any

		assert.Equal(t, want, New(v))
	})

	t.Run("pointer", func(t *testing.T) {
		t.Parallel()

		want := Nullable[*int]{
			value:   nil,
			notNull: true,
		}

		var v *int

		assert.Equal(t, want, New(v))
	})
}

func TestConvertRef(t *testing.T) {
	t.Parallel()

	v := "test"
	tests := []struct {
		name string
		arg  *string
		want Nullable[string]
	}{
		{
			name: "positive",
			arg:  &v,
			want: Nullable[string]{
				value:   v,
				notNull: true,
			},
		},
		{
			name: "nil",
			arg:  nil,
			want: Nullable[string]{},
		},
	}

	conv := func(ref *string) string {
		if ref == nil {
			return ""
		}

		return *ref
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ConvertRef(test.arg, conv)
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestFromRef(t *testing.T) {
	t.Parallel()

	v := 42.42
	tests := []struct {
		name string
		arg  *float64
		want Nullable[float64]
	}{
		{
			name: "nil",
			arg:  nil,
			want: Nullable[float64]{},
		},
		{
			name: "positive",
			arg:  &v,
			want: Nullable[float64]{
				value:   v,
				notNull: true,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FromRef(test.arg)
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestNullable_IsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  Nullable[int]
		want bool
	}{
		{
			name: "empty",
			arg:  Nullable[int]{},
			want: true,
		},
		{
			name: "positive",
			arg: Nullable[int]{
				value:   7,
				notNull: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := test.arg.IsNull()
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestNullable_NotNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  Nullable[string]
		want bool
	}{
		{
			name: "positive",
			arg: Nullable[string]{
				value:   "testing",
				notNull: true,
			},
			want: true,
		},
		{
			name: "empty",
			arg:  Nullable[string]{},
			want: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := test.arg.NotNull()
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestNullable_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  Nullable[float64]
		want float64
	}{
		{
			name: "positive",
			arg: Nullable[float64]{
				value:   13.777,
				notNull: true,
			},
			want: 13.777,
		},
		{
			name: "empty",
			arg:  Nullable[float64]{},
			want: 0,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := test.arg.Value()
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestNullable_Ref(t *testing.T) {
	t.Parallel()

	v := "arg"
	tests := []struct {
		name string
		arg  Nullable[string]
		want *string
	}{
		{
			name: "positive",
			arg: Nullable[string]{
				value:   v,
				notNull: true,
			},
			want: &v,
		},
		{
			name: "empty",
			arg:  Nullable[string]{},
			want: nil,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := test.arg.Ref()
			assert.Equal(t, test.want, got, "incorrect result")
		})
	}
}

func TestNullable_UnmarshalJSON(t *testing.T) {
	type testdata struct {
		Integer Nullable[int]    `json:"integer"`
		String  Nullable[string] `json:"string"`
	}

	type args struct {
		data []byte
	}

	type testCase[T any] struct {
		name    string
		args    args
		want    T
		wantErr assert.ErrorAssertionFunc
	}

	tests := []testCase[testdata]{
		{
			name: "successfully unmarshall empty struct",
			args: args{
				data: []byte(`{}`),
			},
			want: testdata{
				Integer: Nullable[int]{},
				String:  Nullable[string]{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "successfully unmarshall half empty struct",
			args: args{
				data: []byte(`{"string": "kek"}`),
			},
			want: testdata{
				Integer: Nullable[int]{},
				String:  New[string]("kek"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "successfully unmarshall struct",
			args: args{
				data: []byte(`{"string": "kek", "integer": 123}`),
			},
			want: testdata{
				Integer: New[int](123),
				String:  New[string]("kek"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "successfully unmarshall null values struct",
			args: args{
				data: []byte(`{"string": null, "integer": null}`),
			},
			want: testdata{
				Integer: Nullable[int]{},
				String:  Nullable[string]{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "successfully unmarshall wrong struct",
			args: args{
				data: []byte(`{"string": 123, "integer": null}`),
			},
			want: testdata{
				Integer: Nullable[int]{},
				String:  Nullable[string]{},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n testdata

			tt.wantErr(
				t,
				json.Unmarshal(tt.args.data, &n),
				fmt.Sprintf("UnmarshalJSON(%v)", tt.args.data),
			)

			assert.Equal(t, tt.want, n)
		})
	}
}
