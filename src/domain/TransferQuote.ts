import { CurrencyPair } from './CurrencyPair.js';
import { Money } from './Money.js';

export interface TransferQuote {
  pair: CurrencyPair;
  exchangeRate: number;
  amountToSend: Money;
  amountToReceive: Money;
  fee: Money;
  totalDebit: Money;
}
