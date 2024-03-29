import React from 'react';
import PropTypes from 'prop-types';
import style from './index.module.css';

const NUMBER_PATTERN = /^[+-]?\d+\.?\d*$/;

const getValue = (value) => {
  if (!NUMBER_PATTERN.test(value)) {
    return '';
  }

  return value < 0 ? style.red : style.green;
};

export default function PerformanceCell({ type, value }) {
  return (
    <span className={`${style.performance} ${style[type]} ${getValue(value)}`}>
      {value}
    </span>
  );
}

PerformanceCell.propTypes = {
  type: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
};

PerformanceCell.defaultProps = {
  value: '',
};

PerformanceCell.displayName = 'PerformanceCell';
