import { DomainError } from './DomainError.js';

export class UnsupportedCurrencyPairError extends DomainError {
  constructor(sourceCurrency: string, targetCurrency: string) {
    super(`Currency pair ${sourceCurrency}/${targetCurrency} is not supported`);
  }
}
