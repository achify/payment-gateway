import { randomUUID } from 'crypto';
import { PayoutDestination } from '../../domain/PayoutDestination.js';
import { TransferQuote } from '../../domain/TransferQuote.js';
import { PayoutProcessor } from '../../application/ports/PayoutProcessor.js';

export class InMemoryPayoutProcessor implements PayoutProcessor {
  async execute(_request: {
    walletId: string;
    destination: PayoutDestination;
    quote: TransferQuote;
    reference: string;
  }): Promise<{ payoutId: string; status: 'pending' | 'completed' }> {
    return {
      payoutId: randomUUID(),
      status: 'pending'
    };
  }
}
