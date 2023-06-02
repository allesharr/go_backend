package routers

import "testing"

//Little test to be sure in regex
func TestRegex(t *testing.T) {
	data := []string{
		"2023-02-02",
		"2023-13-24",
		"2023-10-33",
		"1999-10-21",
		"21/10/2017",
		"32 Mar 2022",
	}

	for index, elem := range data {
		if index == 0 && !isDataCorrect(elem) {
			t.Error("Correct is not working")
		}
		if index > 0 && isDataCorrect(elem) {
			t.Error("Uncorrect is working")
		}
	}

}
