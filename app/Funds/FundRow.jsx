import React from 'react';
import FaBarChart from 'react-icons/lib/fa/bar-chart';
import FundsService from './FundsService';
import PerformanceCell from './PerformanceCell';
import style from './Funds.css';

function isUpdateTimeGreaterThanSixHours(updateTime) {
  return (Date.now() - Date.parse(updateTime)) / 1000 / 3600 > 6;
}

const MorningStarRow = ({ fund, filterBy }) => (
  <span className={style.row}>
    <span className={style.isin}>{fund.isin}</span>
    <span className={style.label} title={fund.label}>{fund.label}</span>
    <button
      title={fund.category}
      className={style.category}
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
    <a
      title={`Lien vers la source des données. ${fund.ts ? `MàJ le ${new Date(Date.parse(fund.ts)).toString()}` : ''}`}
      href={fund.id && FundsService.getDataUrl(fund.id)}
      rel="noopener noreferrer"
      target="_blank"
      className={`${style.icon} ${fund.ts ? (isUpdateTimeGreaterThanSixHours(fund.ts) ? style.old : style.fresh}) : ''`}
    >
      <FaBarChart />
    </a>
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
