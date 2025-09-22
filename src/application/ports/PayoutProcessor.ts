import { PayoutDestination } from '../../domain/PayoutDestination.js';
import { TransferQuote } from '../../domain/TransferQuote.js';

export interface PayoutProcessor {
  execute(request: {
    walletId: string;
    destination: PayoutDestination;
    quote: TransferQuote;
    reference: string;
  }): Promise<{ payoutId: string; status: 'pending' | 'completed' }>;
}
