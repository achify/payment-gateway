import { describe, expect, it } from 'vitest';
import { TransferCalculationService } from '../TransferCalculationService.js';
import { InMemoryExchangeRateProvider } from '../../infrastructure/exchange/InMemoryExchangeRateProvider.js';
import { TieredFlatFeePolicy } from '../../infrastructure/fees/TieredFlatFeePolicy.js';
import { InvalidAmountError } from '../../domain/errors/InvalidAmountError.js';
import { UnsupportedCurrencyPairError } from '../../domain/errors/UnsupportedCurrencyPairError.js';

const setupService = () => {
  const exchangeRateProvider = new InMemoryExchangeRateProvider({
    'GBP_NGN': 2100
  });
  const feePolicy = new TieredFlatFeePolicy([
    { maxAmountInclusive: 50, fee: 1.99 }
  ], 'GBP');
  const service = new TransferCalculationService(exchangeRateProvider, feePolicy);
  return { service };
};

describe('TransferCalculationService', () => {
  it('calculates quote with fees', async () => {
    const { service } = setupService();

    const quote = await service.calculate({
      sourceCurrency: 'GBP',
      targetCurrency: 'NGN',
      amount: 50
    });

    expect(quote.exchangeRate).toBe(2100);
    expect(quote.amountToReceive.amount).toBe(105000);
    expect(quote.fee.amount).toBe(1.99);
    expect(quote.totalDebit.amount).toBe(51.99);
  });

  it('throws when amount is invalid', async () => {
    const { service } = setupService();

    await expect(
      service.calculate({ sourceCurrency: 'GBP', targetCurrency: 'NGN', amount: 0 })
    ).rejects.toBeInstanceOf(InvalidAmountError);
  });

  it('throws when currency pair unsupported', async () => {
    const { service } = setupService();

    await expect(
      service.calculate({ sourceCurrency: 'USD', targetCurrency: 'NGN', amount: 10 })
    ).rejects.toBeInstanceOf(UnsupportedCurrencyPairError);
  });
});
