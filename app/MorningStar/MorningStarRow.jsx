import React from 'react';
import PerformanceCell from './PerformanceCell';
import style from './MorningStar.css';

const MorningStarRow = ({ performance, filterBy }) => (
  <span className={style.row}>
    <span
      title={`${performance.isin} - ${performance.label}`}
      className={`${style.label} ${style.ellipsis}`}
    >{performance.label}</span>
    <button
      title={performance.category}
      className={`${style.category} ${style.ellipsis}`}
      onClick={() => filterBy && filterBy('category', performance.category)}
    >{performance.category}</button>
    <button
      className={style.rating}
      onClick={() => filterBy && filterBy('rating', performance.rating)}
    >{performance.rating}</button>
    <PerformanceCell value={performance['1m']} type="p1m" />
    <PerformanceCell value={performance['3m']} type="p3m" />
    <PerformanceCell value={performance['6m']} type="p6m" />
    <PerformanceCell value={performance['1y']} type="p1y" />
    <PerformanceCell value={performance.v1y} type="pvol" />
    <PerformanceCell value={performance.score} type="pscore" />
  </span>
);

MorningStarRow.displayName = 'MorningStarRow';

const STRING_OR_NUMBER = React.PropTypes.oneOfType([
  React.PropTypes.string,
  React.PropTypes.number,
]).isRequired;

MorningStarRow.propTypes = {
  performance: React.PropTypes.shape({
    isin: STRING_OR_NUMBER,
    label: STRING_OR_NUMBER,
    category: STRING_OR_NUMBER,
    rating: STRING_OR_NUMBER,
    '1m': STRING_OR_NUMBER,
    '3m': STRING_OR_NUMBER,
    '6m': STRING_OR_NUMBER,
    '1y': STRING_OR_NUMBER,
    v1y: STRING_OR_NUMBER,
    score: STRING_OR_NUMBER,
  }),
  filterBy: React.PropTypes.func,
};

export default MorningStarRow;
