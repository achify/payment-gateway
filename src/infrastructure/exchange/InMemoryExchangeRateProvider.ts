import { CurrencyPair, pairKey } from '../../domain/CurrencyPair.js';
import { ExchangeRateProvider, ExchangeRateQuote } from '../../application/ports/ExchangeRateProvider.js';

export class InMemoryExchangeRateProvider implements ExchangeRateProvider {
  private readonly rates: Map<string, number>;

  constructor(initialRates: Record<string, number>) {
    this.rates = new Map(Object.entries(initialRates));
  }

  async getQuote(pair: CurrencyPair): Promise<ExchangeRateQuote | null> {
    const key = pairKey(pair);
    const rate = this.rates.get(key);
    if (rate === undefined) {
      return null;
    }
    return { pair, rate };
  }

  setRate(pair: CurrencyPair, rate: number): void {
    this.rates.set(pairKey(pair), rate);
  }
}
