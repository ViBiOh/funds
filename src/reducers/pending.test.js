import reducer, { initialState } from './pending';

it('should have a default empty state', () => {
  expect(reducer(undefined, { type: 'ON_CHANGE' })).toEqual(initialState);
});

it('should return given state if action type does not match', () => {
  expect(reducer(initialState, { type: 'ON_CHANGE' })).toEqual(initialState);
});

it('should turn on pending if action type match pattern', () => {
  expect(reducer(initialState, { type: 'ACTION_REQUEST' })).toEqual({
    ...initialState,
    ACTION: true,
  });
});

it('should return given state if pattern match but pending not present', () => {
  expect(reducer(initialState, { type: 'VALIDATE_SUCCEEDED' })).toEqual(initialState);
});

it('should turn off pending if action type match pattern', () => {
  expect(reducer({ ...initialState, ACTION: true }, { type: 'ACTION_SUCCEEDED' })).toEqual({
    ...initialState,
    ACTION: false,
  });
});
