import React from 'react';
import ReactDOM from 'react-dom';
import rollbar from './externals/rollbar';
import Constants from './Constants';
import FundsContainer from './Funds/FundsContainer';

Constants.init().then((config) => {
  if (config.ROLLBAR_TOKEN) {
    rollbar(config.ROLLBAR_TOKEN, config.ENVIRONMENT);
  }

  ReactDOM.render(<FundsContainer />, document.getElementById('root'));
});
