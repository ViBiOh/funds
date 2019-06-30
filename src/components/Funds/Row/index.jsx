import React from 'react';
import PropTypes from 'prop-types';
import { STRING_OR_NUMBER } from 'helpers/Props';
import Button from 'components/Button';
import PerformanceCell from 'components/Funds/PerformanceCell';
import style from './index.module.css';

/**
 * Row functional component
 * @param {Object} options.fund       Fund displayed on th erow
 * @param {Function} options.filterBy Filter function
 */
export default function Row({ fund, filterBy }) {
  return (
    <span className={style.row}>
      <span className={style.isin}>{fund.isin}</span>
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
        <span>{fund.rating}</span>
      </Button>
      <PerformanceCell value={fund['1m']} type="p1m" />
      <PerformanceCell value={fund['3m']} type="p3m" />
      <PerformanceCell value={fund['6m']} type="p6m" />
      <PerformanceCell value={fund['1y']} type="p1y" />
      <PerformanceCell value={fund.v3y} type="pvol" />
      <PerformanceCell value={fund.score} type="pscore" />
    </span>
  );
}

Row.displayName = 'Row';

Row.propTypes = {
  fund: PropTypes.shape({
    isin: PropTypes.string,
    label: PropTypes.string,
    category: PropTypes.string,
    rating: STRING_OR_NUMBER.isRequired,
    '1m': STRING_OR_NUMBER.isRequired,
    '3m': STRING_OR_NUMBER.isRequired,
    '6m': STRING_OR_NUMBER.isRequired,
    '1y': STRING_OR_NUMBER.isRequired,
    v3y: STRING_OR_NUMBER.isRequired,
    score: STRING_OR_NUMBER.isRequired,
  }).isRequired,
  filterBy: PropTypes.func,
};

Row.defaultProps = {
  filterBy: () => null,
};
