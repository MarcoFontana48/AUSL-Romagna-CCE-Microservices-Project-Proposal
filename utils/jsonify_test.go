package utils

import (
	"encoding/json"
	"testing"
)

// test data structures
type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type InvalidStruct struct {
	Ch chan int // channels cannot be marshaled to JSON
}

func TestToJsonByte(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "valid struct",
			input:   TestStruct{Name: "John", Age: 30},
			wantErr: false,
		},
		{
			name:    "valid string",
			input:   "hello world",
			wantErr: false,
		},
		{
			name:    "valid int",
			input:   42,
			wantErr: false,
		},
		{
			name:    "valid slice",
			input:   []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "valid map",
			input:   map[string]int{"a": 1, "b": 2},
			wantErr: false,
		},
		{
			name:    "nil value",
			input:   nil,
			wantErr: false,
		},
		{
			name:    "invalid struct with channel",
			input:   InvalidStruct{Ch: make(chan int)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJsonByte(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJsonByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify the output is valid JSON by unmarshaling it
				var result interface{}
				if err := json.Unmarshal(got, &result); err != nil {
					t.Errorf("ToJsonByte() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestToJsonString(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "valid struct",
			input:   TestStruct{Name: "Alice", Age: 25},
			want:    `{"name":"Alice","age":25}`,
			wantErr: false,
		},
		{
			name:    "valid string",
			input:   "test",
			want:    `"test"`,
			wantErr: false,
		},
		{
			name:    "valid int",
			input:   123,
			want:    "123",
			wantErr: false,
		},
		{
			name:    "valid bool",
			input:   true,
			want:    "true",
			wantErr: false,
		},
		{
			name:    "valid slice",
			input:   []string{"a", "b", "c"},
			want:    `["a","b","c"]`,
			wantErr: false,
		},
		{
			name:    "nil value",
			input:   nil,
			want:    "null",
			wantErr: false,
		},
		{
			name:    "invalid struct with channel",
			input:   InvalidStruct{Ch: make(chan int)},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJsonString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJsonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJsonString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// benchmark tests
func BenchmarkToJsonByte(b *testing.B) {
	data := TestStruct{Name: "BenchmarkTest", Age: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToJsonByte(data)
	}
}

func BenchmarkToJsonString(b *testing.B) {
	data := TestStruct{Name: "BenchmarkTest", Age: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToJsonString(data)
	}
}

// test with complex nested structures
func TestToJsonByteComplexStruct(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Address Address  `json:"address"`
		Hobbies []string `json:"hobbies"`
	}

	person := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
		},
		Hobbies: []string{"reading", "swimming"},
	}

	got, err := ToJsonByte(person)
	if err != nil {
		t.Errorf("ToJsonByte() error = %v", err)
		return
	}

	// Verify it's valid JSON
	var result Person
	if err := json.Unmarshal(got, &result); err != nil {
		t.Errorf("ToJsonByte() produced invalid JSON: %v", err)
	}

	// Verify the data is correct
	if result.Name != person.Name {
		t.Errorf("Name mismatch: got %v, want %v", result.Name, person.Name)
	}
}

// test edge cases
func TestToJsonEdgeCases(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		type Empty struct{}
		empty := Empty{}
		got, err := ToJsonString(empty)
		if err != nil {
			t.Errorf("ToJsonString() error = %v", err)
		}
		if got != "{}" {
			t.Errorf("ToJsonString() = %v, want {}", got)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		var slice []int
		got, err := ToJsonString(slice)
		if err != nil {
			t.Errorf("ToJsonString() error = %v", err)
		}
		if got != "null" {
			t.Errorf("ToJsonString() = %v, want null", got)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := make(map[string]int)
		got, err := ToJsonString(m)
		if err != nil {
			t.Errorf("ToJsonString() error = %v", err)
		}
		if got != "{}" {
			t.Errorf("ToJsonString() = %v, want {}", got)
		}
	})
}
