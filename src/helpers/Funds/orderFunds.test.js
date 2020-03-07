import { orderFunds } from './index';

it('should handle undefined key', () => {
  expect(orderFunds([], {})).toEqual([]);
});

it('should handle ascending order', () => {
  const funds = [
    { score: 4 },
    { score: 3 },
    { score: 2 },
    { score: 5 },
    { score: 1 },
  ];

  const orderedFunds = [
    { score: 1 },
    { score: 2 },
    { score: 3 },
    { score: 4 },
    { score: 5 },
  ];

  expect(
    orderFunds(funds, {
      key: 'score',
    }),
  ).toEqual(orderedFunds);
});

it('should handle descending order', () => {
  const funds = [
    { score: 4 },
    { score: 3 },
    { score: 2 },
    { score: 5 },
    { score: 1 },
  ];

  const orderedFunds = [
    { score: 5 },
    { score: 4 },
    { score: 3 },
    { score: 2 },
    { score: 1 },
  ];

  expect(
    orderFunds(funds, {
      key: 'score',
      descending: true,
    }),
  ).toEqual(orderedFunds);
});

it('should put undefined value at end and handle equality', () => {
  const funds = [
    { score: 4 },
    { score: 4 },
    { isin: 'L0987654321' },
    { score: 3 },
    { score: 2 },
    { score: 5 },
    { score: 1 },
  ];

  const orderedFunds = [
    { score: 5 },
    { score: 4 },
    { score: 4 },
    { score: 3 },
    { score: 2 },
    { score: 1 },
    { isin: 'L0987654321' },
  ];

  expect(
    orderFunds(funds, {
      key: 'score',
      descending: true,
    }),
  ).toEqual(orderedFunds);
});
