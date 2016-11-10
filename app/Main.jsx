import React from 'react';

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
