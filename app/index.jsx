import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, browserHistory } from 'react-router';

import Funds from './Funds/Funds';

ReactDOM.render(
  <Router history={browserHistory}>
    <Route path="/" component={Funds} />
  </Router>,
  document.getElementById('root'),
);
