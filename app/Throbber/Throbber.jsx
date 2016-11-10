import React from 'react';
import style from './throbber.css';

const Throbber = ({ label, title }) => (
  <div className={style.throbberContainer} title={title}>
    {label && <span>{label}</span>}
    <div className={style.throbber}>
      <div className={style.bounce1} />
      <div className={style.bounce2} />
      <div className={style.bounce3} />
    </div>
  </div>
);

Throbber.propTypes = {
  label: React.PropTypes.string,
  title: React.PropTypes.string,
};

Throbber.displayname = 'Throbber';

export default Throbber;
