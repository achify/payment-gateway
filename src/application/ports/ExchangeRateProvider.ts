import { CurrencyPair } from '../../domain/CurrencyPair.js';

export interface ExchangeRateQuote {
  pair: CurrencyPair;
  rate: number;
}

export interface ExchangeRateProvider {
  getQuote(pair: CurrencyPair): Promise<ExchangeRateQuote | null>;
}
