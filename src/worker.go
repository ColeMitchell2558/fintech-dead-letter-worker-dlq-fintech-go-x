package src

import "encoding/json"

type PaymentEvent struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"`
	Attempts int    `json:"attempts"`
}
type Decision struct{ Notify, Escalate bool }

func Decide(e PaymentEvent) Decision {
	return Decision{Notify: true, Escalate: e.Amount >= 100000 || e.Attempts >= 3}
}
func Handle(c *Client, m Message) (Decision, error) {
	var e PaymentEvent
	if err := json.Unmarshal([]byte(m.Payload), &e); err != nil {
		return Decision{}, err
	}
	d := Decide(e)
	if err := c.QueueAck(m.MessageID); err != nil {
		return d, err
	}
	return d, nil
}
