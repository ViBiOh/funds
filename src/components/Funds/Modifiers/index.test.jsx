import { fireEvent, render } from '@testing-library/react';
import React from 'react';
import Modifiers from './index';

function defaultProps() {
  return {
    fundsSize: 0,
    initialSize: 0,
    filters: {},
    filterBy: () => null,
    order: {},
    orderBy: () => null,
    reverseOrder: () => null,
    aggregat: {},
    aggregateBy: () => null,
    onAggregateSizeChange: () => null,
  };
}

it('should always render as a div', () => {
  const props = defaultProps();
  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelector('div')).toBeTruthy();
});

it('should not render count if size match', () => {
  const props = defaultProps();
  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-count]').length).toEqual(0);
});

it('should not render filters if empty', () => {
  const props = defaultProps();
  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-filters]').length).toEqual(0);
});

it('should not render order if not specified', () => {
  const props = defaultProps();
  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-order]').length).toEqual(0);
});

it('should render count if size differ', () => {
  const props = defaultProps();
  props.fundsSize = 4000;
  props.initialSize = 8000;

  const { container } = render(<Modifiers {...props} />);
  const fundsCountElements = container.querySelectorAll('[data-funds-count]');
  expect(fundsCountElements.length).toEqual(1);
  expect(fundsCountElements[0].textContent).toEqual('4000 / 8000');
});

it('should render each filter separately', () => {
  const props = defaultProps();
  props.filters = {
    isin: '123456',
    category: 'Test',
    score: 0,
  };

  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-filter]').length).toEqual(2);
});

it('should call given callback for removing filter', () => {
  const props = defaultProps();
  props.filters = {
    isin: '123456',
  };
  props.filterBy = jest.fn();

  const { container } = render(<Modifiers {...props} />);

  const clearButton = container.querySelector(
    'span[data-funds-filter] button[data-funds-filter-clear]',
  );
  fireEvent.click(clearButton);

  expect(props.filterBy).toHaveBeenCalledWith('isin', '');
});

it('should render order if specified', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
  };

  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-order]').length).toEqual(1);
  expect(container.querySelectorAll('[data-fa-sort-amount-up]').length).toEqual(
    1,
  );
});

it('should render order descending if specified', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
    descending: true,
  };

  const { container } = render(<Modifiers {...props} />);
  expect(
    container.querySelectorAll('[data-fa-sort-amount-down]').length,
  ).toEqual(1);
});

it('should call given callback for removing order', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
  };
  props.orderBy = jest.fn();

  const { container } = render(<Modifiers {...props} />);
  const button = container.querySelector(
    '[data-funds-order] button[data-funds-order-clear]',
  );
  fireEvent.click(button);

  expect(props.orderBy).toHaveBeenCalledWith('');
});

it('should render aggregat if specified', () => {
  const props = defaultProps();
  props.aggregat = {
    key: 'score',
  };

  const { container } = render(<Modifiers {...props} />);
  expect(container.querySelectorAll('[data-funds-aggregat]').length).toEqual(1);
});

it('should call given callback for removing order', () => {
  const props = defaultProps();
  props.aggregat = {
    key: 'score',
  };
  props.aggregateBy = jest.fn();

  const { container } = render(<Modifiers {...props} />);
  const button = container.querySelector(
    '[data-funds-aggregat] button[data-funds-aggregat-clear]',
  );
  fireEvent.click(button);

  expect(props.aggregateBy).toHaveBeenCalledWith('');
});
