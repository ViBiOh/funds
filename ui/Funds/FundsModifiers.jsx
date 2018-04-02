import React from 'react';
import PropTypes from 'prop-types';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaPieChart from 'react-icons/lib/fa/pie-chart';
import Button from '../Button/Button';
import { COLUMNS, AGGREGATE_SIZES } from './FundsConstantes';
import style from './FundsModifiers.less';

function renderCount(fundsSize, initialSize) {
  if (fundsSize === initialSize) {
    return null;
  }

  return (
    <span key="count" className={style.modifier}>
      {fundsSize} / {initialSize}
    </span>
  );
}

function renderFilters(filters, filterBy) {
  return Object.keys(filters)
    .filter(filter => filters[filter])
    .map(filter => (
      <span key={filter} className={style.modifier}>
        <FaFilter />
        <span>
          <em> {COLUMNS[filter].label}</em> ≃{' '}
        </span>
        <span>{filters[filter]}</span>
        <Button type="none" onClick={() => filterBy(filter, '')}>
          <FaClose />
        </Button>
      </span>
    ));
}

function renderOrder(order, orderBy, reverseOrder) {
  return (
    order.key && (
      <span className={style.modifier}>
        <Button type="none" onClick={reverseOrder}>
          {order.descending ? <FaSortAmountDesc /> : <FaSortAmountAsc />}
        </Button>
        <span>{COLUMNS[order.key].label}</span>
        <Button type="none" onClick={() => orderBy('')}>
          <FaClose />
        </Button>
      </span>
    )
  );
}

function renderAggregat(aggregat, aggregateBy, onAggregateSizeChange) {
  if (!aggregat.key) {
    return null;
  }

  const { label } = COLUMNS[aggregat.key];

  return (
    <span className={style.modifier}>
      <FaPieChart />
      <select value={aggregat.size} onChange={onAggregateSizeChange}>
        {AGGREGATE_SIZES.map(size => (
          <option key={size} value={size}>
            {size}
          </option>
        ))}
      </select>{' '}
      {label}
      <Button type="none" onClick={() => aggregateBy('')}>
        <FaClose />
      </Button>
    </span>
  );
}

const FundsModifier = ({
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
}) => (
  <div className={style.list}>
    {renderCount(fundsSize, initialSize)}
    {renderFilters(filters, filterBy)}
    {renderOrder(order, orderBy, reverseOrder)}
    {renderAggregat(aggregat, aggregateBy, onAggregateSizeChange)}
  </div>
);

FundsModifier.displayName = 'FundsModifier';

FundsModifier.propTypes = {
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

export default FundsModifier;
