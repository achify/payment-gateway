import { Money } from '../../domain/Money.js';

export interface FeePolicy {
  calculateFee(amount: Money): Money;
}
