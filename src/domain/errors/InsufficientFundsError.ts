import { DomainError } from './DomainError.js';

export class InsufficientFundsError extends DomainError {
  constructor(balance: number, required: number) {
    super(`Insufficient funds. Balance: ${balance}, required: ${required}`);
  }
}
