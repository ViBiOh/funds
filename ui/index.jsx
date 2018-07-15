import React from 'react';
import ReactDOM from 'react-dom';
import Fathom from './external/fathom';
import Constants from './Constants';
import FundsContainer from './Funds/FundsContainer';

Constants.init().then((config) => {
  ReactDOM.render(<FundsContainer />, document.getElementById('root'));
  Fathom.init(config);
});
