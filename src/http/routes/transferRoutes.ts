import { Router } from 'express';
import { z } from 'zod';
import { TransferCalculationService } from '../../application/TransferCalculationService.js';
import { mapErrorToHttp } from '../HttpErrorMapper.js';

export const createTransferRoutes = (calculationService: TransferCalculationService): Router => {
  const router = Router();

  const requestSchema = z.object({
    sourceCurrency: z.string().min(3),
    targetCurrency: z.string().min(3),
    amount: z.number().positive()
  });

  router.post('/calculate', async (req, res) => {
    try {
      const body = requestSchema.parse(req.body);
      const quote = await calculationService.calculate(body);
      res.status(200).json({
        exchangeRate: quote.exchangeRate,
        amountToSend: quote.amountToSend,
        amountToReceive: quote.amountToReceive,
        fee: quote.fee,
        totalDebit: quote.totalDebit
      });
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(422).json({ message: error.errors.map((e) => e.message).join(', ') });
        return;
      }
      const httpError = mapErrorToHttp(error);
      res.status(httpError.status).json(httpError.body);
    }
  });

  return router;
};
