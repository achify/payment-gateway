import { InsufficientFundsError } from './errors/InsufficientFundsError.js';
import { InvalidAmountError } from './errors/InvalidAmountError.js';

export interface WalletProps {
  id: string;
  currency: string;
  balance: number;
}

export class Wallet {
  readonly id: string;
  readonly currency: string;
  private _balance: number;

  constructor({ id, currency, balance }: WalletProps) {
    this.id = id;
    this.currency = currency;
    this._balance = balance;
  }

  get balance(): number {
    return this._balance;
  }

  debit(amount: number): void {
    if (amount <= 0) {
      throw new InvalidAmountError(amount);
    }
    if (this._balance < amount) {
      throw new InsufficientFundsError(this._balance, amount);
    }
    this._balance = Number((this._balance - amount).toFixed(2));
  }
}
