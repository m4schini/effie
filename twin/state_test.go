package match

import "testing"

func TestMatchStateMachine(t *testing.T) {
	m := New("sum1D")

	if m.state != startState {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnFoundData(nil)
	if m.state != Started {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnFoundData(nil)
	if m.state != InGame {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnFoundData(nil)
	if m.state != InGame {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnNoData()
	if m.state != Stopped {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnNoData()
	if m.state != Idle {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnNoData()
	if m.state != Idle {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnFoundData(nil)
	if m.state != Started {
		t.Log("machine is in unexpected state")
		t.Fail()
	}

	m.OnNoData()
	if m.state != Stopped {
		t.Log("machine is in unexpected state")
		t.Fail()
	}
}
