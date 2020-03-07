import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import style from './index.module.css';

/**
 * Button.
 * @param Object} props Props of the component.
 * @return {React.Component} Button with rendered children.
 */
export default function Button({
  children,
  type,
  active,
  className,
  ...buttonProps
}) {
  let content = children;
  if (Array.isArray(children)) {
    content = <div className={style.wrapper}>{children}</div>;
  }

  const classes = classnames(style.button, style[type], className, {
    [style.active]: active,
  });

  return (
    <button type="button" className={classes} {...buttonProps}>
      {content}
    </button>
  );
}

Button.displayName = 'Button';

Button.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
  type: PropTypes.oneOf([
    'transparent',
    'primary',
    'success',
    'info',
    'warning',
    'danger',
    'none',
  ]),
  active: PropTypes.bool,
  className: PropTypes.string,
};

Button.defaultProps = {
  children: null,
  type: 'primary',
  active: false,
  className: '',
};
