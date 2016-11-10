import React from 'react';
import { Link } from 'react-router';

const toggleProps = {
  idle: <span>&#x2715;</span>,
  active: <span>&#x2261;</span>,
};

const Main = ({ children }) => (
  <span id="mainLayout">
    <header>
      <h1>Funds</h1>
    </header>
    <article>{children}</article>
  </span>
);

Main.propTypes = {
  children: React.PropTypes.element.isRequired,
};

export default Main;
