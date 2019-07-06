import actions from 'actions';
import reducer, { initialState } from './funds';

it('should return initial state', () => {
  expect(reducer(undefined, { type: '' })).toEqual(initialState);
});

it('should update with given funds on fetch succeed', () => {
  const all = [{ isin: 8000 }];

  expect(reducer(initialState, actions.getFundsSucceeded(all))).toEqual({
    ...initialState,
    all,
    displayed: all,
  });
});

it('should remove funds without id', () => {
  const all = [{ isin: 8000 }, { name: 'test' }];

  expect(reducer(initialState, actions.getFundsSucceeded(all))).toEqual({
    ...initialState,
    all: [{ isin: 8000 }],
    displayed: [{ isin: 8000 }],
  });
});

it('should add given filter', () => {
  expect(reducer(initialState, actions.setFilter('name', 'test'))).toEqual({
    ...initialState,
    filters: {
      name: 'test',
    },
  });
});

it('should concat filters', () => {
  let state = reducer(
    { ...initialState, filters: { ...initialState.filters, previous: 'next' } },
    actions.setFilter('name', 'test'),
  );

  state = reducer(state, actions.setFilter('previous', true));

  expect(state).toEqual({
    ...initialState,
    filters: {
      name: 'test',
      previous: true,
    },
  });
});

it('should apply filters', () => {
  const all = [{ isin: 8, name: 'Vite' }, { isin: 1000, name: 'Emile' }];

  expect(
    reducer({ ...initialState, all, displayed: all }, actions.setFilter('name', 'ite')),
  ).toEqual({
    ...initialState,
    all,
    displayed: [{ isin: 8, name: 'Vite' }],
    filters: {
      name: 'ite',
    },
  });
});
