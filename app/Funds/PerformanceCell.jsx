import React from 'react';
import PropTypes from 'prop-types';
import style from './PerformanceCell.less';

const NUMBER_PATTERN = /^[+-]?[0-9]+\.?[0-9]*$/;

const getValue = (value) => {
  if (!NUMBER_PATTERN.test(value)) {
    return '';
  }

  return value < 0 ? style.red : style.green;
};

const PerformanceCell = ({ type, value }) => (
  <span className={`${style.performance} ${style[type]} ${getValue(value)}`}>{value}</span>
);

PerformanceCell.propTypes = {
  type: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
};

PerformanceCell.defaultProps = {
  value: '',
};

PerformanceCell.displayName = 'PerformanceCell';

export default PerformanceCell;
