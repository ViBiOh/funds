import { fireEvent, render } from '@testing-library/react';
import React from 'react';
import Row from './index';

function defaultProps() {
  return {
    fund: {
      category: 'Test',
      rating: 4,
      '1m': 0.0,
      '3m': 0.0,
      '6m': 0.0,
      '1y': 0.0,
      v3y: 0.0,
      score: 0.0,
    },
  };
}

it('should always render as a span', () => {
  const props = defaultProps();
  const { queryAllByRole, container } = render(<Row {...props} />);

  const button = queryAllByRole('button')[0];
  fireEvent.click(button);

  expect(container.querySelector('span')).toBeTruthy();
});

it('should call given filterBy func on category click', () => {
  const props = defaultProps();
  props.filterBy = jest.fn();
  const { queryAllByRole } = render(<Row {...props} />);

  const button = queryAllByRole('button')[0];
  fireEvent.click(button);

  expect(props.filterBy).toHaveBeenCalledWith('category', 'Test');
});

it('should call given filterBy func on rating click', () => {
  const props = defaultProps();
  props.filterBy = jest.fn();
  const { queryAllByRole } = render(<Row {...props} />);

  const button = queryAllByRole('button')[1];
  fireEvent.click(button);

  expect(props.filterBy).toHaveBeenCalledWith('rating', 4);
});
