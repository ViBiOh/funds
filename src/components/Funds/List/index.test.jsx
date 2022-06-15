import { render } from '@testing-library/react';
import React from 'react';
import List from './index';

function defaultProps() {
  return {
    funds: [],
    filterBy: () => null,
  };
}

it('should always render as a div', () => {
  const props = defaultProps();
  const { container } = render(<List {...props} />);
  expect(container.querySelector('div')).toBeTruthy();
});

it('should render at least one header row', () => {
  const props = defaultProps();
  const { queryByTestId } = render(<List {...props} />);
  expect(queryByTestId('row-header')).toBeTruthy();
});

it('should render one row per funds', () => {
  const props = defaultProps();
  props.funds = [
    {
      id: '1000',
      rating: 1,
      '1m': 0.0,
      '3m': 0.0,
      '6m': 0.0,
      '1y': 0.0,
      v3y: 0.0,
      score: 0.0,
    },
    {
      id: '2000',
      rating: 1,
      '1m': 0.0,
      '3m': 0.0,
      '6m': 0.0,
      '1y': 0.0,
      v3y: 0.0,
      score: 0.0,
    },
    {
      id: '8000',
      rating: 1,
      '1m': 0.0,
      '3m': 0.0,
      '6m': 0.0,
      '1y': 0.0,
      v3y: 0.0,
      score: 0.0,
    },
  ];
  const { queryAllByTestId } = render(<List {...props} />);

  expect(queryAllByTestId('fund-row').length).toEqual(3);
});
