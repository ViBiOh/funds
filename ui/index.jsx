import React from 'react';
import ReactDOM from 'react-dom';
import Constants from './Constants';
import FundsContainer from './Funds/FundsContainer';

Constants.init().then(() => {
  ReactDOM.render(<FundsContainer />, document.getElementById('root'));
});
