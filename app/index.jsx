import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';

import Main from './Main';
import MorningStarList from './MorningStar/MorningStar';

ReactDOM.render(
  <Router history={browserHistory}>
    <Route path="/" component={Main}>
      <IndexRoute component={MorningStar} />
    </Route>
  </Router>,
  document.getElementById('root'),
);
