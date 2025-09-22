import { FeePolicy } from '../../application/ports/FeePolicy.js';
import { Money } from '../../domain/Money.js';

export interface FlatFeeTier {
  maxAmountInclusive: number;
  fee: number;
}

export class TieredFlatFeePolicy implements FeePolicy {
  constructor(private readonly tiers: FlatFeeTier[], private readonly currency: string) {}

  calculateFee(amount: Money): Money {
    if (amount.currency !== this.currency) {
      return { amount: 0, currency: amount.currency };
    }

    const tier = this.tiers.find((t) => amount.amount <= t.maxAmountInclusive);
    if (!tier) {
      return { amount: 0, currency: amount.currency };
    }

    return {
      amount: Number(tier.fee.toFixed(2)),
      currency: amount.currency
    };
  }
}
