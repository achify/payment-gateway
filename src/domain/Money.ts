export interface Money {
  amount: number;
  currency: string;
}

export const formatMoney = (money: Money): string => {
  return `${money.currency} ${money.amount.toFixed(2)}`;
};
