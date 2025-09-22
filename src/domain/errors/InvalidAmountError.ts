import { DomainError } from './DomainError.js';

export class InvalidAmountError extends DomainError {
  constructor(amount: number) {
    super(`Amount must be greater than zero. Received: ${amount}`);
  }
}
