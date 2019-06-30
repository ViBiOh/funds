import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import Constants from 'Constants';
import App from 'App';
import appStore from 'Store';
import './index.css';

(async () => {
  await Constants.init();
  ReactDOM.render(
    <Provider store={appStore}>
      <App />
    </Provider>,
    document.getElementById('root'),
  );
})();
