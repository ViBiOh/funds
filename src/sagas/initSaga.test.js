import { call, put } from 'redux-saga/effects';
import Config from 'services/Config';
import actions from 'actions';
import initSaga from './initSaga';

it('should call Config.getConfig', () => {
  const iterator = initSaga();

  expect(iterator.next().value).toEqual(call(Config.getConfig));
});

it('should put success after API call', () => {
  const iterator = initSaga();
  iterator.next();

  const value = { url: 'localhost' };

  expect(iterator.next(value).value).toEqual(put(actions.getConfigSucceeded(value)));
});

it('should put error on failure', () => {
  const iterator = initSaga();
  iterator.next();

  const value = new Error('Test');

  expect(iterator.throw(value).value).toEqual(put(actions.getConfigFailed(value)));
});
