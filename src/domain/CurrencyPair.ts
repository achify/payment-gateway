export interface CurrencyPair {
  source: string;
  target: string;
}

export const normalizeCurrency = (currency: string): string => currency.trim().toUpperCase();

export const createCurrencyPair = (source: string, target: string): CurrencyPair => ({
  source: normalizeCurrency(source),
  target: normalizeCurrency(target)
});

export const pairKey = (pair: CurrencyPair): string => `${pair.source}_${pair.target}`;
