package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {

	name := "Chris"
	age := 20
	age2 := 39
	city := "Amsterdam"
	city2 := "London"
	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct {
				Name string
			}{name},
			[]string{name},
		},
		{
			"Struct with no string fields",
			struct {
				Age int
			}{age},
			[]string{},
		},
		{
			"Struct with two string fields",
			struct {
				Name string
				City string
			}{name, city},
			[]string{name, city},
		},
		{
			"Struct with nested fields",
			Person{
				name,
				Profile{age, city},
			},
			[]string{name, city},
		},
		{
			"Struct with pointers",
			&Person{
				Name:    name,
				Profile: Profile{age, city},
			},
			[]string{name, city},
		},
		{
			"Slices",
			[]Profile{
				{age, city},
				{age2, city2},
			},
			[]string{city, city2},
		},
		{
			"Arrays",
			[2]Profile{
				{age, city},
				{age2, city2},
			},
			[]string{city, city2},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := make([]string, 0)
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("Maps", func(t *testing.T) {
		aMap := map[string]string{
			"Foo": "Bar",
			"Baz": "Boz",
		}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Bar")
		assertContains(t, got, "Boz")
	})
}

func assertContains(t *testing.T, haystack []string, needle string) {
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
			return
		}
	}
	if !contains {
		t.Errorf("Expected %+v to contain %q, but it didn't", haystack, needle)
	}
}
