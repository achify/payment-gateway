import { Router } from 'express';
import { z } from 'zod';
import { PayoutService } from '../../application/PayoutService.js';
import { mapErrorToHttp } from '../HttpErrorMapper.js';

export const createPayoutRoutes = (payoutService: PayoutService): Router => {
  const router = Router();

  const destinationSchema = z.object({
    type: z.literal('bank_account'),
    bank: z.object({
      country: z.string().min(2),
      bankCode: z.string().min(2),
      accountNumber: z.string().min(5),
      accountName: z.string().min(3)
    })
  });

  const requestSchema = z.object({
    walletId: z.string().min(1),
    amount: z.number().positive(),
    sourceCurrency: z.string().min(3),
    targetCurrency: z.string().min(3),
    reference: z.string().min(1),
    destination: destinationSchema
  });

  router.post('/', async (req, res) => {
    try {
      const body = requestSchema.parse(req.body);
      const result = await payoutService.initiatePayout(body);
      res.status(202).json(result);
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
