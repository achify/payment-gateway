import { describe, expect, it, vi } from 'vitest';
import { PayoutService } from '../PayoutService.js';
import { InMemoryWalletRepository } from '../../infrastructure/repositories/InMemoryWalletRepository.js';
import { InMemoryExchangeRateProvider } from '../../infrastructure/exchange/InMemoryExchangeRateProvider.js';
import { TieredFlatFeePolicy } from '../../infrastructure/fees/TieredFlatFeePolicy.js';
import { PayoutProcessor } from '../ports/PayoutProcessor.js';
import { EventPublisher } from '../ports/EventPublisher.js';
import { Wallet } from '../../domain/Wallet.js';
import { InsufficientFundsError } from '../../domain/errors/InsufficientFundsError.js';
import { UnsupportedCurrencyPairError } from '../../domain/errors/UnsupportedCurrencyPairError.js';

describe('PayoutService', () => {
  const setup = () => {
    const walletRepository = new InMemoryWalletRepository([
      new Wallet({ id: 'wallet-1', currency: 'GBP', balance: 100 })
    ]);
    const exchangeRateProvider = new InMemoryExchangeRateProvider({ 'GBP_NGN': 2100 });
    const feePolicy = new TieredFlatFeePolicy([{ maxAmountInclusive: 50, fee: 1.99 }], 'GBP');
    const payoutProcessor: PayoutProcessor = {
      execute: vi.fn().mockResolvedValue({ payoutId: 'payout-123', status: 'pending' })
    };
    const eventPublisher: EventPublisher = {
      publish: vi.fn().mockResolvedValue(undefined)
    };

    const payoutService = new PayoutService(
      walletRepository,
      exchangeRateProvider,
      feePolicy,
      payoutProcessor,
      eventPublisher
    );

    return { payoutService, payoutProcessor, eventPublisher, walletRepository };
  };

  it('initiates payout when funds sufficient', async () => {
    const { payoutService, payoutProcessor, eventPublisher, walletRepository } = setup();

    const response = await payoutService.initiatePayout({
      walletId: 'wallet-1',
      amount: 50,
      sourceCurrency: 'GBP',
      targetCurrency: 'NGN',
      reference: 'INV-001',
      destination: {
        type: 'bank_account',
        bank: {
          country: 'NG',
          bankCode: '044',
          accountNumber: '1234567890',
          accountName: 'John Doe'
        }
      }
    });

    expect(response.payoutId).toBe('payout-123');
    expect(response.quote.totalDebit.amount).toBe(51.99);
    expect(response.walletBalance).toBeCloseTo(48.01, 2);
    expect(payoutProcessor.execute).toHaveBeenCalledOnce();
    expect(eventPublisher.publish).toHaveBeenCalledWith('payouts.initiated', expect.any(Object));

    const savedWallet = await walletRepository.findById('wallet-1');
    expect(savedWallet?.balance).toBeCloseTo(48.01, 2);
  });

  it('fails when funds insufficient', async () => {
    const { payoutService } = setup();

    await expect(
      payoutService.initiatePayout({
        walletId: 'wallet-1',
        amount: 120,
        sourceCurrency: 'GBP',
        targetCurrency: 'NGN',
        reference: 'INV-002',
        destination: {
          type: 'bank_account',
          bank: {
            country: 'NG',
            bankCode: '044',
            accountNumber: '1234567890',
            accountName: 'John Doe'
          }
        }
      })
    ).rejects.toBeInstanceOf(InsufficientFundsError);
  });

  it('fails when wallet currency mismatch', async () => {
    const { payoutService } = setup();

    await expect(
      payoutService.initiatePayout({
        walletId: 'wallet-1',
        amount: 10,
        sourceCurrency: 'USD',
        targetCurrency: 'NGN',
        reference: 'INV-003',
        destination: {
          type: 'bank_account',
          bank: {
            country: 'NG',
            bankCode: '044',
            accountNumber: '1234567890',
            accountName: 'John Doe'
          }
        }
      })
    ).rejects.toBeInstanceOf(UnsupportedCurrencyPairError);
  });
});
