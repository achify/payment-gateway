import { createCurrencyPair } from '../domain/CurrencyPair.js';
import { Money } from '../domain/Money.js';
import { TransferQuote } from '../domain/TransferQuote.js';
import { InvalidAmountError } from '../domain/errors/InvalidAmountError.js';
import { UnsupportedCurrencyPairError } from '../domain/errors/UnsupportedCurrencyPairError.js';
import { ExchangeRateProvider } from './ports/ExchangeRateProvider.js';
import { FeePolicy } from './ports/FeePolicy.js';

export interface TransferCalculationRequest {
  sourceCurrency: string;
  targetCurrency: string;
  amount: number;
}

export class TransferCalculationService {
  constructor(
    private readonly exchangeRateProvider: ExchangeRateProvider,
    private readonly feePolicy: FeePolicy
  ) {}

  async calculate(request: TransferCalculationRequest): Promise<TransferQuote> {
    if (request.amount <= 0) {
      throw new InvalidAmountError(request.amount);
    }

    const pair = createCurrencyPair(request.sourceCurrency, request.targetCurrency);
    const quote = await this.exchangeRateProvider.getQuote(pair);
    if (!quote) {
      throw new UnsupportedCurrencyPairError(pair.source, pair.target);
    }

    const amountToSend: Money = {
      amount: Number(request.amount.toFixed(2)),
      currency: pair.source
    };
    const amountToReceive: Money = {
      amount: Number((amountToSend.amount * quote.rate).toFixed(2)),
      currency: pair.target
    };
    const fee = this.feePolicy.calculateFee(amountToSend);
    const totalDebit: Money = {
      amount: Number((amountToSend.amount + fee.amount).toFixed(2)),
      currency: pair.source
    };

    return {
      pair,
      exchangeRate: quote.rate,
      amountToSend,
      amountToReceive,
      fee,
      totalDebit
    };
  }
}
