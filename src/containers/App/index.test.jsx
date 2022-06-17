import React from 'react';
import { render } from '@testing-library/react';

import { AppComponent } from './index';

jest.mock('../Funds', () => () => {
  const MockName = 'connected-funds';
  return <MockName />;
});

function defaultProps() {
  return {
    ready: false,
    init: () => null,
  };
}

it('should render as a div if not ready', () => {
  const props = defaultProps();
  const { container } = render(<AppComponent {...props} />);
  expect(container.querySelector('div')).toBeTruthy();
});

it('should show a Throbber while not ready', () => {
  const props = defaultProps();

  const { queryByText } = render(<AppComponent {...props} />);
  expect(queryByText('Loading environment')).toBeTruthy();
});

it('should trigger init on mount', () => {
  const props = defaultProps();
  props.init = jest.fn();

  render(<AppComponent {...props} />);

  expect(props.init).toHaveBeenCalled();
});

it('should render Funds when ready', () => {
  const props = defaultProps();
  props.ready = true;

  const { container } = render(<AppComponent {...props} />);

  expect(container.querySelector('connected-funds')).toBeTruthy();
});
