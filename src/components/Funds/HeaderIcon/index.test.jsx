import React from 'react';
import { COLUMNS } from 'components/Funds/Constants';
import { render, fireEvent } from '@testing-library/react';
import HeaderIcon from './index';

function defaultProps() {
  return {
    filter: 'sortable',
    onClick: () => null,
    icon: <span />,
    displayed: false,
  };
}

it('should always render as a span', () => {
  const props = defaultProps();
  const { container } = render(<HeaderIcon {...props} />);
  expect(container.querySelector('span')).toBeTruthy();
});

it('should call given callback on item click', () => {
  const props = defaultProps();
  props.onClick = jest.fn();

  const { queryAllByRole } = render(<HeaderIcon {...props} />);

  fireEvent.click(queryAllByRole('button')[0]);

  expect(props.onClick).toHaveBeenCalledWith(Object.keys(COLUMNS)[0]);
});
