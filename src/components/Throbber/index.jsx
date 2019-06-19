import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import style from './index.module.css';

/**
 * Throbber for displaying background task.
 * @param {Object} props Props of the component.
 * @return {React.Component} Throbber with label and title if provided
 */
export default function Throbber({ label, title, className }) {
  const classes = classnames(style.throbber, className);

  return (
    <div className={style.container} title={title}>
      {label && <span>{label}</span>}
      <div className={classes}>
        <div className={style.bounce1} />
        <div className={style.bounce2} />
        <div className={style.bounce3} />
      </div>
    </div>
  );
}

Throbber.displayname = 'Throbber';

Throbber.propTypes = {
  label: PropTypes.string,
  title: PropTypes.string,
  className: PropTypes.string,
};

Throbber.defaultProps = {
  label: '',
  title: '',
  className: '',
};
