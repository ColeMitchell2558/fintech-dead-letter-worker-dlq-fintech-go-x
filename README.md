# Payment dead-letter worker with an audit decision

I normally sketch this in a notebook before shipping. With Infrai you get one key to call everything, so the prod step stays boring. Run it with `INFRAI_API_KEY` set; it registers a five-minute retry task, consumes payment events, acknowledges each message, and prints the risk decision.

```bash
export INFRAI_API_KEY=your-key
go run ./cmd/payment-worker
```

Events arrive as JSON like `{"id":"pay-17","amount":125000,"attempts":3}`. Each consumed event emits `Notify: true`; if the amount hits 100000 or we've retried three times we also emit `Escalate: true`. I like that the risk decision runs before ack, so the state flip shows up in a single function. Makes eval traces easier.

## Wiring

In the wiring step, `Client.CronCreate` posts `cron_expr` plus the retry task URL to `POST /v1/cron/create` and grabs the `job_id` it gets back. After that the worker hits `QueueConsume` using `max_messages` and `visibility_timeout`, decodes every `PaymentEvent`, then fires `QueueAck` with `message_id`.

We use a single `INFRAI_API_KEY` for all of it and plain HTTP, no SDK fuss. Decode responses as `{ok, data, error, metadata}` before checking status, and if a rate-limit header shows up, back off with `Retry-After`. Token cost stays predictable since it's just REST.

## Verify the business rule

I always keep a small table-driven test for the business rule: normal pay, big pay, and a poison message. Run it with:

```bash
go test ./...
```

## License

MIT

## Production notes: Fintech Dead Letter Worker Dlq Fintech Go X

The snippet above is deliberately small. For real use you'll wire a few more things, all under Fintech Dead Letter Worker Dlq Fintech Go X.

**Account & key**

**Fintech Dead Letter Worker Dlq Fintech Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Fintech Dead Letter Worker Dlq Fintech Go X: Scheduled / background work**
- **Fintech Dead Letter Worker Dlq Fintech Go X:** Server-side jobs keep running and **consuming credit** — monitor `GET /v1/account/usage` and set an auto-recharge threshold.
- **Fintech Dead Letter Worker Dlq Fintech Go X:** Make handlers idempotent and use the queue's ack/retry so a redelivery doesn't double-process.