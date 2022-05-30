package cache

import "testing"

func TestGetHungerGamesCount(t *testing.T) {
	count, setCount := GetHungerGamesCount()

	before := count()
	t.Log("count", before)

	setCount(before + 1)

	after := count()
	t.Log("count", after)

	if after != (before + 1) {
		t.Fail()
	}
}
