package src

import "testing"

func TestDecide(t *testing.T) {
	cases := []struct {
		name             string
		e                PaymentEvent
		notify, escalate bool
	}{{"ordinary", PaymentEvent{Amount: 500, Attempts: 1}, true, false}, {"large", PaymentEvent{Amount: 100000, Attempts: 1}, true, true}, {"poison", PaymentEvent{Amount: 500, Attempts: 3}, true, true}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Decide(tc.e)
			if got.Notify != tc.notify || got.Escalate != tc.escalate {
				t.Fatalf("got %+v", got)
			}
		})
	}
}
