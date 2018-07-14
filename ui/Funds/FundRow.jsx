import React from 'react';
import PropTypes from 'prop-types';
import Button from '../Button';
import PerformanceCell from './PerformanceCell';
import style from './FundRow.css';

const FundRow = ({ fund, filterBy }) => (
  <span className={style.row}>
    <span className={style.isin}>
      {fund.isin}
    </span>
    <span className={style.label} title={fund.label}>
      {fund.label}
    </span>
    <Button
      className={style.category}
      type="none"
      onClick={() => filterBy('category', fund.category)}
    >
      <span className={style.ellipsis} title={fund.category}>
        {fund.category}
      </span>
    </Button>
    <Button className={style.rating} type="none" onClick={() => filterBy('rating', fund.rating)}>
      <span>
        {fund.rating}
      </span>
    </Button>
    <PerformanceCell value={fund['1m']} type="p1m" />
    <PerformanceCell value={fund['3m']} type="p3m" />
    <PerformanceCell value={fund['6m']} type="p6m" />
    <PerformanceCell value={fund['1y']} type="p1y" />
    <PerformanceCell value={fund.v3y} type="pvol" />
    <PerformanceCell value={fund.score} type="pscore" />
  </span>
);

FundRow.displayName = 'FundRow';

const STRING_OR_NUMBER = PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired;

FundRow.propTypes = {
  fund: PropTypes.shape({
    isin: PropTypes.string,
    label: PropTypes.string,
    category: PropTypes.string,
    rating: STRING_OR_NUMBER,
    '1m': STRING_OR_NUMBER,
    '3m': STRING_OR_NUMBER,
    '6m': STRING_OR_NUMBER,
    '1y': STRING_OR_NUMBER,
    v3y: STRING_OR_NUMBER,
    score: STRING_OR_NUMBER,
  }).isRequired,
  filterBy: PropTypes.func,
};

FundRow.defaultProps = {
  filterBy: () => null,
};

export default FundRow;
