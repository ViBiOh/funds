import { call, put } from 'redux-saga/effects';
import Funds from 'services/Funds';
import actions from 'actions';
import getFundsSaga from './getFundsSaga';

it('should call Funds.getFunds', () => {
  const iterator = getFundsSaga();

  expect(iterator.next().value).toEqual(call(Funds.getFunds));
});

it('should put success after API call', () => {
  const iterator = getFundsSaga();
  iterator.next();

  const value = [{ id: 8000 }];

  expect(iterator.next(value).value).toEqual(put(actions.getFundsSucceeded(value)));
});

it('should put error on failure', () => {
  const iterator = getFundsSaga();
  iterator.next();

  const value = new Error('Test');

  expect(iterator.throw(value).value).toEqual(put(actions.getFundsFailed(value)));
});
