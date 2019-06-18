import React from 'react';
import ReactDOM from 'react-dom';
import Constants from 'Constants';
import App from 'App';
import './index.css';

Constants.init().then((config) => {
  ReactDOM.render(<App />, document.getElementById('root'));
});
