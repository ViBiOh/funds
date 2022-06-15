import { render } from '@testing-library/react';
import React from 'react';
import PerformanceCell from './index';

const defaultProps = {
  type: '1y',
};

it('should always render as a span', () => {
  const { container } = render(<PerformanceCell {...defaultProps} />);
  expect(container.querySelector('span')).toBeTruthy();
});

it('should have style for given value', () => {
  const { container } = render(
    <PerformanceCell {...defaultProps} value={10} />,
  );
  expect(container.firstChild.className.split(' ').length).toEqual(3);
});

it('should have style for negative value', () => {
  const { container } = render(
    <PerformanceCell {...defaultProps} value={-10} />,
  );

  expect(container.firstChild.className.split(' ').length).toEqual(3);
});
