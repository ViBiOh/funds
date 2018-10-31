import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import rollbar from './externals/rollbar';
import Constants from './Constants';
import App from './App';
import * as serviceWorker from './serviceWorker';

Constants.init().then((config) => {
  if (config.ROLLBAR_TOKEN) {
    rollbar(config.ROLLBAR_TOKEN, config.ENVIRONMENT);
  }

  ReactDOM.render(<App />, document.getElementById('root'));
});

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
