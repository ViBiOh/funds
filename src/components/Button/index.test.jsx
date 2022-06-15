import React from 'react';
import { render } from '@testing-library/react';
import Button from './index';

it('should always render as a button', () => {
  const { queryByRole } = render(<Button />);
  expect(queryByRole('button')).toBeTruthy();
});

it('should not wrap child', () => {
  const { container } = render(
    <Button>
      <span>First</span>
    </Button>,
  );

  expect(container.querySelector('span').textContent).toEqual('First');
});

it('should wrap children in div', () => {
  const { container } = render(
    <Button>
      <span>First</span>
      <span>Second</span>
    </Button>,
  );

  expect(container.querySelector('div')).toBeTruthy();
});
