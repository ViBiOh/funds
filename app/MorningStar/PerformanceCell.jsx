import React from 'react';
import style from './MorningStar.css';

const NUMBER_PATTERN = /^[+-]?[0-9]+\.?[0-9]*$/;

function getPerformanceStyle(performance) {
  if (!NUMBER_PATTERN.test(performance)) {
    return null;
  }

  return {
    color: performance < 0 ? '#d43f3a' : '#4cae4c',
  };
}

const PerformanceCell = ({ type, value }) => (
  <span
    style={getPerformanceStyle(value)}
    className={`${style.performance} ${style[type]}`}
  >{value}</span>
);

PerformanceCell.propTypes = {
  type: React.PropTypes.string.isRequired,
  value: React.PropTypes.oneOfType([
    React.PropTypes.string,
    React.PropTypes.number,
  ]),
};

PerformanceCell.displayName = 'PerformanceCell';

export default PerformanceCell;
