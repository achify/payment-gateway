import { Wallet } from '../../domain/Wallet.js';
import { WalletRepository } from '../../application/ports/WalletRepository.js';

export class InMemoryWalletRepository implements WalletRepository {
  private readonly wallets: Map<string, Wallet> = new Map();

  constructor(initialWallets: Wallet[] = []) {
    initialWallets.forEach((wallet) => {
      this.wallets.set(wallet.id, wallet);
    });
  }

  async findById(id: string): Promise<Wallet | null> {
    const wallet = this.wallets.get(id);
    if (!wallet) {
      return null;
    }
    return new Wallet({ id: wallet.id, currency: wallet.currency, balance: wallet.balance });
  }

  async save(wallet: Wallet): Promise<void> {
    this.wallets.set(wallet.id, wallet);
  }
}
