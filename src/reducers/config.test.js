import actions from 'actions';
import reducer, { initialState } from './config';

it('should return initial state', () => {
  expect(reducer(undefined, { type: '' })).toEqual(initialState);
});

it('should shotre given config', () => {
  expect(reducer(initialState, actions.getConfigSucceeded({ url: 'localhost' }))).toEqual({
    ready: true,
    url: 'localhost',
  });
});
