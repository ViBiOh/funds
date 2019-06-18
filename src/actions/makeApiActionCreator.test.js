import { makeApiActionCreator } from './creator';

it('should return action type', () => {
  const apiActions = makeApiActionCreator('fetch', ['payload'], ['response']);

  expect(apiActions.FETCH).toEqual('FETCH');
  expect(apiActions.FETCH_REQUEST).toEqual('FETCH_REQUEST');
  expect(apiActions.FETCH_SUCCEEDED).toEqual('FETCH_SUCCEEDED');
  expect(apiActions.FETCH_FAILED).toEqual('FETCH_FAILED');
});

it('should return action creator', () => {
  const apiActions = makeApiActionCreator('fetch', ['payload'], ['response']);

  expect(apiActions.fetch('id')).toEqual({ type: 'FETCH_REQUEST', payload: 'id' });
  expect(apiActions.fetchSucceeded('valid')).toEqual({
    type: 'FETCH_SUCCEEDED',
    response: 'valid',
  });
  expect(apiActions.fetchFailed(new Error('hi'))).toEqual({
    type: 'FETCH_FAILED',
    error: new Error('hi'),
  });
});
