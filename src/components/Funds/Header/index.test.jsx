import React from 'react';
import { render } from '@testing-library/react';
import Header from './index';

function defaultProps() {
  return {
    aggregateBy: () => null,
    filterBy: () => null,
    orderBy: () => null,
  };
}

it('should always render as a header', () => {
  const props = defaultProps();
  const { container } = render(<Header {...props} />);
  expect(container.querySelector('header')).toBeTruthy();
});
