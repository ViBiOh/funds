import { render } from '@testing-library/react';
import React from 'react';
import Graph from './index';

jest.mock('chart.js');

function defaultProps() {
  return {
    aggregat: {},
    aggregated: [],
  };
}

it('should not render if no aggregat key', () => {
  const props = defaultProps();
  const { container } = render(<Graph {...props} />);
  expect(container.firstChild).toBeNull();
});

it('should render a Graph if key and aggregated values', () => {
  const props = defaultProps();
  props.aggregat.key = 'test';
  props.aggregated = [
    {
      label: 'first',
      count: 8,
    },
    {
      label: 'second',
      count: 12,
    },
  ];

  const { queryByTestId } = render(<Graph {...props} />);
  expect(queryByTestId('graph')).toBeTruthy();
});
