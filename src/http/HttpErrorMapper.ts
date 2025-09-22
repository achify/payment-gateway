import { DomainError } from '../domain/errors/DomainError.js';
import { InsufficientFundsError } from '../domain/errors/InsufficientFundsError.js';
import { InvalidAmountError } from '../domain/errors/InvalidAmountError.js';
import { UnsupportedCurrencyPairError } from '../domain/errors/UnsupportedCurrencyPairError.js';

export const mapErrorToHttp = (error: unknown): { status: number; body: { message: string } } => {
  if (error instanceof UnsupportedCurrencyPairError) {
    return { status: 400, body: { message: error.message } };
  }

  if (error instanceof InvalidAmountError) {
    return { status: 422, body: { message: error.message } };
  }

  if (error instanceof InsufficientFundsError) {
    return { status: 409, body: { message: error.message } };
  }

  if (error instanceof DomainError) {
    return { status: 400, body: { message: error.message } };
  }

  return { status: 500, body: { message: 'Internal Server Error' } };
};
