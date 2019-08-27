package reflection

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {

	type Profile struct {
		Age  int
		City string
	}
	tests := []struct {
		desc     string
		input    interface{}
		expected []string
	}{
		{"Struct with one string field", struct {
			Name string
		}{"Joe Sixpack"}, []string{"Joe Sixpack"}},

		{"Struct with two string field", struct {
			Name string
			City string
		}{"Jane Doe", "Chicago"}, []string{"Jane Doe", "Chicago"}},

		{"Struct with a non string field", struct {
			Name string
			Age  int
		}{"Hal", 42}, []string{"Hal"}},

		{"Nested fields", struct {
			Name    string
			Profile Profile
		}{"Lara Croft", Profile{28, "London"}}, []string{"Lara Croft", "London"}},

		{"With Pointers", &struct {
			Name    string
			Profile Profile
		}{"Eric Cartman", Profile{28, "Southpark"}}, []string{"Eric Cartman", "Southpark"}},

		{"With Slices", []Profile{{44, "D.C."}, {56, "Miami"}}, []string{"D.C.", "Miami"}},

		{"With Arrays", [2]Profile{{14, "Langhorne"}, {8, "Duluth"}}, []string{"Langhorne", "Duluth"}},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got []string
			Walk(tt.input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}

	t.Run("With Maps", func(t *testing.T) {
		var got []string
		Walk(map[string]string{"Foo": "Bar", "Baz": "Boz"}, func(input string) {
			got = append(got, input)
		})

		for _, want := range []string{"Bar", "Boz"} {
			assertContains(t, got, want)
		}
	})
}

func assertContains(t *testing.T, haystack []string, needle string) {
	t.Helper()

	contains := false

	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %+v to contain %q, but it did not", haystack, needle)
	}
}
