import { render } from '@testing-library/react';
import React from 'react';
import Throbber from './index';

it('should render into a div', () => {
  const { container } = render(<Throbber />);
  expect(container.querySelector('div')).toBeTruthy();
});

it('should have no label by default', () => {
  const { container } = render(<Throbber />);
  expect(container.querySelectorAll('span').length).toEqual(0);
});

it('should have label when given', () => {
  const { queryByText } = render(<Throbber label="test" />);
  expect(queryByText('test')).toBeTruthy();
});
