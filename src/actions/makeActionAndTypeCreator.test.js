import { makeActionAndTypeCreator } from './creator';

it('should return action type', () => {
  expect(
    makeActionAndTypeCreator('ACTION_TYPE', 'actionType').ACTION_TYPE,
  ).toEqual('ACTION_TYPE');
});

it('should return action creator', () => {
  expect(
    makeActionAndTypeCreator('ACTION_TYPE', 'actionType', [
      'payload',
    ]).actionType('content'),
  ).toEqual({
    type: 'ACTION_TYPE',
    payload: 'content',
  });
});
