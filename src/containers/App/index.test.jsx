import React from 'react';
import funtch from 'funtch';
import PropTypes from 'prop-types';
import { act, render } from '@testing-library/react';
import { Provider } from 'react-redux';
import { combineReducers, createStore } from 'redux';
import error, { initialState as errorInitialState } from 'reducers/error';
import funds, { initialState as fundsInitialState } from 'reducers/funds';
import pending, { initialState as pendingInitialState } from 'reducers/pending';
import config, { initialState as configInitialState } from 'reducers/config';

import { AppComponent } from './index';

jest.mock('funtch');
jest.mock('AppStore');

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

function ReduxProvider({ children }) {
  return (
    <Provider
      store={createStore(
        combineReducers({
          config,
          error,
          funds,
          pending,
        }),
        {
          config: configInitialState,
          error: errorInitialState,
          funds: fundsInitialState,
          pending: pendingInitialState,
        },
      )}
    >
      {children}
    </Provider>
  );
}

ReduxProvider.propTypes = { children: PropTypes.node.isRequired };

it('should render Funds when ready', async (done) => {
  const props = defaultProps();

  funtch.get.mockResolvedValue({});

  const { container } = render(<AppComponent {...props} />, {
    wrapper: ReduxProvider,
  });

  const connectedFunds = await container.querySelectorAll(
    '[data-connected-funds]',
  );

  expect(connectedFunds.length).toEqual(1);
  done();
});
