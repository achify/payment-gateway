# Payment Gateway API (Go)

This project models a payment gateway that can quote foreign exchange transfers and initiate payouts while keeping business
logic isolated from delivery and infrastructure details. The implementation is written in Go and demonstrates SOLID-friendly
composition, clean layering, and strong reliance on interfaces so that databases, Kafka, and other integrations can be
exchanged without touching the domain.

## Architecture

```
cmd/server              -> application wiring & HTTP bootstrap
internal/
  application/          -> use-case services and DTOs
    ports/              -> interfaces for infrastructure dependencies
  domain/               -> entities, value objects, and domain errors
  http/                 -> HTTP handlers that adapt requests/responses
  infrastructure/       -> in-memory adapters for exchange rates, fees, wallets, payouts, and events
```

Key principles:

- **Segregation of concerns** – each layer has a single responsibility; handlers only translate transport concerns, services
  implement business rules, and infrastructure contains integrations.
- **Dependency inversion** – application services depend only on interfaces declared under `internal/application/ports`.
- **Testability** – in-memory adapters make it straightforward to unit test success/failure paths and to swap in fakes during
  integration or end-to-end tests.

## Running the service

```bash
go run ./cmd/server
```

The server listens on `http://localhost:8080`.

### Transfer quote

`POST /transfer/quote`

```json
{
  "sourceCurrency": "GBP",
  "targetCurrency": "NGN",
  "amount": "50"
}
```

Response

```json
{
  "canTransact": true,
  "exchangeRate": "2100",
  "sendAmount": { "currency": "GBP", "amount": "50.00" },
  "receiveAmount": { "currency": "NGN", "amount": "105000.00" },
  "fee": { "currency": "GBP", "amount": "1.99" },
  "totalDebit": { "currency": "GBP", "amount": "51.99" }
}
```

### Initiate payout

`POST /payouts`

```json
{
  "walletId": "wallet-123",
  "currency": "GBP",
  "amount": "50",
  "reference": "client-ref-1",
  "destination": {
    "type": "bank_account",
    "accountNumber": "1234567890",
    "bankCode": "999",
    "bankName": "Demo Bank",
    "country": "NG",
    "currency": "NGN",
    "recipientName": "John Doe"
  }
}
```

Response

```json
{
  "payoutId": "payout-000001",
  "quote": { /* same structure as transfer quote */ },
  "remainingBalance": { "currency": "GBP", "amount": "48.01" }
}
```

## Testing

```bash
go test ./...
```

Unit tests cover:

- successful transfer quote calculation with fee application
- payout initiation happy path (debit, processor call, event emission)
- failures for unsupported currency pairs, processor errors, and insufficient wallet funds

These tests are inexpensive and can be extended into E2E tests by swapping the in-memory adapters for HTTP/Kafka fakes.

## Infrastructure touchpoints & trade-offs

- **Database bottleneck** – the `WalletRepository` interface is the single place that would hit a persistent data store. The
  in-memory implementation mimics optimistic locking by copying structs; a real implementation should use transactions or row
  locks to avoid double spending.
- **Kafka integration** – events flow through the `EventPublisher` interface. Replacing the logging publisher with a Kafka
  producer centralises the change to one adapter while keeping business logic untouched. Back-pressure and delivery semantics
  (at-least vs. exactly-once) are the main trade-offs to evaluate.
- **Monetary precision** – the domain uses `shopspring/decimal` to avoid floating point rounding. Production systems should
  persist amounts in minor units and enforce currency-specific precision.
- **Idempotency & retries** – payout orchestration assumes a single successful call. Real gateways add idempotency keys, retry
  policies, and saga orchestration when dealing with external processors and message brokers.

## Extending the system

- Swap `InMemoryWalletRepository` for a PostgreSQL implementation that satisfies `WalletRepository`.
- Replace `LoggingEventPublisher` with a Kafka-backed publisher to stream payout lifecycle events.
- Implement additional fee policies (percentage tiers, partner overrides) by creating new structs that satisfy `FeePolicy`.
- Add authentication, request tracing, and observability middlewares without touching the application layer.
