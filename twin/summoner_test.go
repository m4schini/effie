package twin

import (
	"context"
	"testing"
	"time"
)

func TestTwin(t *testing.T) {

	st := NewSummonerTwin("PDSOL-mJByhFpQOMv-7KrpKpS9sutwhg0TIyfdNxtIhgcx9CvWlZKbidYw")

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
	st.AutoRefresh(ctx)

	if ctx.Err() != nil {
		t.Fail()
	}
}

func TestMod(t *testing.T) {
	numbers := []int{
		123,
		124,
		5,
		456,
		867,
		1,
		234,
		2,
	}

	for _, number := range numbers {
		t.Log(number % 10)
	}
}
