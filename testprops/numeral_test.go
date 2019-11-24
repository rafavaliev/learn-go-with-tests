package testprops

import "testing"

func TestRomanNumerals(t *testing.T) {

	cases := []struct {
		Description string
		Arabic      int
		Want        string
	}{
		{"1 to I", 1, "I"},
		{"2 to II", 2, "II"},
		{"3 to III", 3, "III"},
		{"4 to IV", 4, "IV"},
		{"5 to V", 5, "V"},
		{"9 to IX", 9, "IX"},
		{"20 to XX", 20, "XX"},
		{"39 to XXXIX", 39, "XXXIX"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			if got != test.Want {
				t.Errorf("got %q, want %q", got, test.Want)
			}
		})
	}

}
