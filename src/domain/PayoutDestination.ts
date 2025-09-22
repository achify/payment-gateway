export interface BankAccountDetails {
  country: string;
  bankCode: string;
  accountNumber: string;
  accountName: string;
}

export interface PayoutDestination {
  type: 'bank_account';
  bank: BankAccountDetails;
}
