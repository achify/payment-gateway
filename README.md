# Payment Gateway API

This project demonstrates a layered payment gateway API that separates domain logic from infrastructure concerns. It exposes endpoints for quoting cross-border transfers and initiating payouts while keeping dependencies behind interfaces to enable testing and future integrations.

## Architecture Overview

The codebase follows SOLID principles and a clean architecture inspired structure:

- **Domain layer (`src/domain`)** – value objects, entities, and domain-specific errors.
- **Application layer (`src/application`)** – orchestrates use cases through services and works only with interfaces defined under `ports`.
- **Infrastructure layer (`src/infrastructure`)** – adapters that implement the ports (e.g., in-memory exchange rates, wallet repository, payout processor, event publisher, fee policy).
- **HTTP layer (`src/http`)** – Express routes that translate HTTP requests/responses into use-case invocations.

The services depend on interfaces (`ExchangeRateProvider`, `WalletRepository`, `PayoutProcessor`, `EventPublisher`, `FeePolicy`), making it straightforward to swap implementations for actual databases or messaging systems.

## Endpoints

| Method | Path | Description |
| ------ | ---- | ----------- |
| `POST` | `/api/transfers/calculate` | Returns an FX quote, fee, and wallet debit total. |
| `POST` | `/api/payouts` | Initiates a payout using a previously calculated quote. |

### Transfer calculation request

```json
{
  "sourceCurrency": "GBP",
  "targetCurrency": "NGN",
  "amount": 50
}
```

### Transfer calculation response

```json
{
  "exchangeRate": 2100,
  "amountToSend": { "amount": 50, "currency": "GBP" },
  "amountToReceive": { "amount": 105000, "currency": "NGN" },
  "fee": { "amount": 1.99, "currency": "GBP" },
  "totalDebit": { "amount": 51.99, "currency": "GBP" }
}
```

### Payout request

```json
{
  "walletId": "wallet-gbp-1",
  "amount": 50,
  "sourceCurrency": "GBP",
  "targetCurrency": "NGN",
  "reference": "INV-2024-05",
  "destination": {
    "type": "bank_account",
    "bank": {
      "country": "NG",
      "bankCode": "044",
      "accountNumber": "1234567890",
      "accountName": "John Doe"
    }
  }
}
```

### Payout response

```json
{
  "payoutId": "...",
  "status": "pending",
  "quote": { /* same structure as the calculation response */ },
  "walletBalance": 48.01
}
```

## Running the project

```bash
npm install
npm run dev
```

The API will be available on `http://localhost:3000`.

## Testing

Unit tests are provided for both transfer calculations and payout initiation:

```bash
npm test
```

The tests cover success paths and failure scenarios such as unsupported currency pairs, invalid amounts, and insufficient funds.

## Replacing infrastructure dependencies

- **Database access** – swap the `InMemoryWalletRepository` implementation for one that reads/writes to a persistent store (e.g., PostgreSQL). All interactions with the database are isolated within the repository.
- **Kafka publishing** – replace `LoggingEventPublisher` with a Kafka-backed publisher. `PayoutService` only depends on the `EventPublisher` interface, so the change is limited to wiring.

## Design trade-offs & bottlenecks

- **Precision** – For brevity, amounts are handled as JavaScript numbers. In production, use a decimal library to avoid rounding issues.
- **Idempotency** – The payout endpoint assumes single invocation. Real systems should add idempotency keys and retry-safe processor implementations.
- **Concurrency** – The in-memory wallet repository is single-threaded. A real database implementation would need proper locking/transactions.
- **External integrations** – Stubbing the payout processor and publisher keeps the example simple while showing where real HTTP/Kafka clients would plug in.

