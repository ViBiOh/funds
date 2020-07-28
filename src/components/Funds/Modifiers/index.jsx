import React from 'react';
import PropTypes from 'prop-types';
import {
  FaWindowClose,
  FaFilter,
  FaSortAmountUp,
  FaSortAmountDown,
  FaChartPie,
} from 'react-icons/fa';
import Button from 'components/Button';
import { COLUMNS, AGGREGATE_SIZES } from 'components/Funds/Constants';
import style from './index.module.css';

function renderCount(fundsSize, initialSize) {
  if (fundsSize === initialSize) {
    return null;
  }

  return (
    <span key="count" className={style.modifier} data-funds-count>
      {fundsSize}
      {' / '}
      {initialSize}
    </span>
  );
}

function renderFilters(filters, filterBy) {
  return Object.keys(filters)
    .filter((filter) => filters[filter])
    .map((filter) => (
      <span key={filter} className={style.modifier} data-funds-filter>
        <FaFilter />
        <span>
          <em>
            {' '}
            {COLUMNS[filter].label}
          </em>
          {' '}
          ≃
          {' '}
        </span>

        <span>{filters[filter]}</span>
        <Button
          type="none"
          onClick={() => filterBy(filter, '')}
          data-funds-filter-clear
        >
          <FaWindowClose />
        </Button>
      </span>
    ));
}

function renderOrder(order, orderBy, reverseOrder) {
  if (!order.key) {
    return null;
  }

  return (
    <span className={style.modifier} data-funds-order>
      <Button type="none" onClick={reverseOrder}>
        {order.descending ? <FaSortAmountDown /> : <FaSortAmountUp />}
      </Button>

      <span>{COLUMNS[order.key].label}</span>
      <Button type="none" onClick={() => orderBy('')} data-funds-order-clear>
        <FaWindowClose />
      </Button>
    </span>
  );
}

function renderAggregat(aggregat, aggregateBy, onAggregateSizeChange) {
  if (!aggregat.key) {
    return null;
  }

  const { label } = COLUMNS[aggregat.key];

  return (
    <span className={style.modifier} data-funds-aggregat>
      <FaChartPie />
      <select value={aggregat.size} onChange={onAggregateSizeChange}>
        {AGGREGATE_SIZES.map((size) => (
          <option key={size} value={size}>
            {size}
          </option>
        ))}
      </select>
      {' '}
      {label}
      <Button
        type="none"
        onClick={() => aggregateBy('')}
        data-funds-aggregat-clear
      >
        <FaWindowClose />
      </Button>
    </span>
  );
}

export default function Modifiers({
  fundsSize,
  initialSize,
  filters,
  filterBy,
  order,
  orderBy,
  reverseOrder,
  aggregat,
  aggregateBy,
  onAggregateSizeChange,
}) {
  return (
    <div className={style.list}>
      {renderCount(fundsSize, initialSize)}
      {renderFilters(filters, filterBy)}
      {renderOrder(order, orderBy, reverseOrder)}
      {renderAggregat(aggregat, aggregateBy, onAggregateSizeChange)}
    </div>
  );
}

Modifiers.displayName = 'Modifiers';

Modifiers.propTypes = {
  fundsSize: PropTypes.number.isRequired,
  initialSize: PropTypes.number.isRequired,
  filters: PropTypes.shape({}).isRequired,
  filterBy: PropTypes.func.isRequired,
  order: PropTypes.shape({}).isRequired,
  orderBy: PropTypes.func.isRequired,
  reverseOrder: PropTypes.func.isRequired,
  aggregat: PropTypes.shape({}).isRequired,
  aggregateBy: PropTypes.func.isRequired,
  onAggregateSizeChange: PropTypes.func.isRequired,
};
