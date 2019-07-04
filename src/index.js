import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import actions from 'actions';
import Constants from 'Constants';
import store from 'AppStore';
import App from 'containers/App';
import './index.css';

store.dispatch(actions.init());

(async () => {
  await Constants.init();
  ReactDOM.render(
    <Provider store={store}>
      <App />
    </Provider>,
    document.getElementById('root'),
  );
})();
