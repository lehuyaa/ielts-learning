package foundation

import "testing"

func TestHelpers(t *testing.T) {
	if sentences("My name is Huy. I am 32! Do you like tea? ok") != 4 {
		t.Fatal("sentences")
	}
	if percent(3, 0) != 0 || percent(5, 4) != 100 || percent(1, 3) != 33 {
		t.Fatal("percent")
	}
	if average(10, 0) != 0 || average(150, 2) != 75 {
		t.Fatal("average")
	}
}
