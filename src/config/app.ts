import express from 'express';
import bodyParser from 'body-parser';
import { TransferCalculationService } from '../application/TransferCalculationService.js';
import { PayoutService } from '../application/PayoutService.js';
import { InMemoryExchangeRateProvider } from '../infrastructure/exchange/InMemoryExchangeRateProvider.js';
import { InMemoryWalletRepository } from '../infrastructure/repositories/InMemoryWalletRepository.js';
import { TieredFlatFeePolicy } from '../infrastructure/fees/TieredFlatFeePolicy.js';
import { InMemoryPayoutProcessor } from '../infrastructure/payouts/InMemoryPayoutProcessor.js';
import { LoggingEventPublisher } from '../infrastructure/events/LoggingEventPublisher.js';
import { createTransferRoutes } from '../http/routes/transferRoutes.js';
import { createPayoutRoutes } from '../http/routes/payoutRoutes.js';
import { Wallet } from '../domain/Wallet.js';

export const createApp = () => {
  const app = express();
  app.use(bodyParser.json());

  const exchangeRateProvider = new InMemoryExchangeRateProvider({
    'GBP_NGN': 2100
  });

  const feePolicy = new TieredFlatFeePolicy(
    [
      { maxAmountInclusive: 50, fee: 1.99 }
    ],
    'GBP'
  );

  const walletRepository = new InMemoryWalletRepository([
    new Wallet({ id: 'wallet-gbp-1', currency: 'GBP', balance: 100 })
  ]);

  const payoutProcessor = new InMemoryPayoutProcessor();
  const eventPublisher = new LoggingEventPublisher(console);

  const calculationService = new TransferCalculationService(exchangeRateProvider, feePolicy);
  const payoutService = new PayoutService(
    walletRepository,
    exchangeRateProvider,
    feePolicy,
    payoutProcessor,
    eventPublisher
  );

  app.use('/api/transfers', createTransferRoutes(calculationService));
  app.use('/api/payouts', createPayoutRoutes(payoutService));

  return app;
};
