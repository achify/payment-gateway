import { createCurrencyPair } from '../domain/CurrencyPair.js';
import { PayoutDestination } from '../domain/PayoutDestination.js';
import { TransferQuote } from '../domain/TransferQuote.js';
import { Wallet } from '../domain/Wallet.js';
import { InsufficientFundsError } from '../domain/errors/InsufficientFundsError.js';
import { UnsupportedCurrencyPairError } from '../domain/errors/UnsupportedCurrencyPairError.js';
import { EventPublisher } from './ports/EventPublisher.js';
import { ExchangeRateProvider } from './ports/ExchangeRateProvider.js';
import { FeePolicy } from './ports/FeePolicy.js';
import { PayoutProcessor } from './ports/PayoutProcessor.js';
import { WalletRepository } from './ports/WalletRepository.js';
import { TransferCalculationService } from './TransferCalculationService.js';

export interface InitiatePayoutRequest {
  walletId: string;
  amount: number;
  sourceCurrency: string;
  targetCurrency: string;
  destination: PayoutDestination;
  reference: string;
}

export interface InitiatePayoutResponse {
  payoutId: string;
  status: 'pending' | 'completed';
  quote: TransferQuote;
  walletBalance: number;
}

export class PayoutService {
  private readonly calculationService: TransferCalculationService;

  constructor(
    private readonly walletRepository: WalletRepository,
    exchangeRateProvider: ExchangeRateProvider,
    feePolicy: FeePolicy,
    private readonly payoutProcessor: PayoutProcessor,
    private readonly eventPublisher: EventPublisher
  ) {
    this.calculationService = new TransferCalculationService(exchangeRateProvider, feePolicy);
  }

  async initiatePayout(request: InitiatePayoutRequest): Promise<InitiatePayoutResponse> {
    const wallet = await this.getWalletOrThrow(request.walletId);
    const pair = createCurrencyPair(request.sourceCurrency, request.targetCurrency);
    if (wallet.currency !== pair.source) {
      throw new UnsupportedCurrencyPairError(wallet.currency, pair.target);
    }

    const quote = await this.calculationService.calculate({
      sourceCurrency: pair.source,
      targetCurrency: pair.target,
      amount: request.amount
    });

    if (wallet.balance < quote.totalDebit.amount) {
      throw new InsufficientFundsError(wallet.balance, quote.totalDebit.amount);
    }

    wallet.debit(quote.totalDebit.amount);
    await this.walletRepository.save(wallet);

    const payoutResult = await this.payoutProcessor.execute({
      walletId: wallet.id,
      destination: request.destination,
      quote,
      reference: request.reference
    });

    await this.eventPublisher.publish('payouts.initiated', {
      payoutId: payoutResult.payoutId,
      walletId: wallet.id,
      reference: request.reference,
      amount: quote.amountToSend,
      fee: quote.fee,
      destination: request.destination,
      status: payoutResult.status
    });

    return {
      payoutId: payoutResult.payoutId,
      status: payoutResult.status,
      quote,
      walletBalance: wallet.balance
    };
  }

  private async getWalletOrThrow(walletId: string): Promise<Wallet> {
    const wallet = await this.walletRepository.findById(walletId);
    if (!wallet) {
      throw new Error(`Wallet ${walletId} not found`);
    }
    return wallet;
  }
}
