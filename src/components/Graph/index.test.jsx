import { render } from '@testing-library/react';
import React from 'react';
import Graph from './index';

jest.mock('chart.js');

function defaultProps() {
  return {
    data: {
      datasets: [],
      labels: [],
    },
  };
}

it('should always render as a canvas', () => {
  const props = defaultProps();
  const { container } = render(<Graph {...props} />);
  expect(container.querySelector('canvas')).toBeTruthy();
});

it('should do nothing if chart not here on unmount', () => {
  const props = defaultProps();
  const { unmount } = render(<Graph {...props} />);
  unmount();
});
