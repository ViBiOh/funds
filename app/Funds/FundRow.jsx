import React from 'react';
import PerformanceCell from './PerformanceCell';
import style from './Funds.css';

const MorningStarRow = ({ fund, filterBy }) => (
  <span className={style.row}>
    <span
      title={`${fund.isin} - ${fund.label}`}
      className={`${style.label} ${style.ellipsis}`}
    >{fund.label}</span>
    <button
      title={fund.category}
      className={`${style.category} ${style.ellipsis}`}
      onClick={() => filterBy && filterBy('category', fund.category)}
    >{fund.category}</button>
    <button
      className={style.rating}
      onClick={() => filterBy && filterBy('rating', fund.rating)}
    >{fund.rating}</button>
    <PerformanceCell value={fund['1m']} type="p1m" />
    <PerformanceCell value={fund['3m']} type="p3m" />
    <PerformanceCell value={fund['6m']} type="p6m" />
    <PerformanceCell value={fund['1y']} type="p1y" />
    <PerformanceCell value={fund.v3y} type="pvol" />
    <PerformanceCell value={fund.score} type="pscore" />
  </span>
);

MorningStarRow.displayName = 'MorningStarRow';

const STRING_OR_NUMBER = React.PropTypes.oneOfType([
  React.PropTypes.string,
  React.PropTypes.number,
]).isRequired;

MorningStarRow.propTypes = {
  fund: React.PropTypes.shape({
    isin: STRING_OR_NUMBER,
    label: STRING_OR_NUMBER,
    category: STRING_OR_NUMBER,
    rating: STRING_OR_NUMBER,
    '1m': STRING_OR_NUMBER,
    '3m': STRING_OR_NUMBER,
    '6m': STRING_OR_NUMBER,
    '1y': STRING_OR_NUMBER,
    v3y: STRING_OR_NUMBER,
    score: STRING_OR_NUMBER,
  }),
  filterBy: React.PropTypes.func,
};

export default MorningStarRow;
