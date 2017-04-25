import React from 'react';
import PropTypes from 'prop-types';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaPieChart from 'react-icons/lib/fa/pie-chart';
import { COLUMNS, CHART_COLORS, AGGREGATE_SIZES } from './FundsConstantes';
import Graph from './Graph';
import style from './FundsModifier.css';

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
  return Object.keys(filters).filter(filter => filters[filter]).map(filter => (
    <span key={filter} className={style.modifier}>
      <span>
        <FaFilter />
      </span>
      <span><em> {COLUMNS[filter].label}</em> ≃ </span>
      {filters[filter]}
      <button onClick={() => filterBy(filter, '')}>
        <FaClose />
      </button>
    </span>
  ));
}

function renderOrder(order, orderBy, reverseOrder) {
  return (
    order.key &&
    <span className={style.modifier}>
      <button onClick={reverseOrder}>
        {order.descending ? <FaSortAmountDesc /> : <FaSortAmountAsc />}
      </button>
      &nbsp;{COLUMNS[order.key].label}
      <button onClick={() => orderBy('')}>
        <FaClose />
      </button>
    </span>
  );
}

function renderAggregate(aggregate, aggregateBy, aggregated, onAggregateSizeChange) {
  if (!aggregate.key) {
    return null;
  }

  const label = COLUMNS[aggregate.key].label;

  const options = {
    legend: false,
    scales: {
      xAxes: [
        {
          display: false,
        },
      ],
      yAxes: [
        {
          display: false,
          ticks: {
            beginAtZero: true,
          },
        },
      ],
    },
  };

  const data = {
    labels: [],
    datasets: [
      {
        label: 'Count',
        data: [],
        backgroundColor: [],
      },
    ],
  };

  let i = 0;
  aggregated.forEach(entry => {
    data.labels.push(entry.label);
    data.datasets[0].data.push(entry.count);
    data.datasets[0].backgroundColor.push(CHART_COLORS[i]);

    i = (i + 1) % CHART_COLORS.length;
  });

  return [
    <span key="label" className={style.modifier}>
      <FaPieChart />
      &nbsp;
      <select value={aggregate.size} onChange={onAggregateSizeChange}>
        {AGGREGATE_SIZES.map(size => <option key={size} value={size}>{size}</option>)}
      </select> | {label}
      <button onClick={() => aggregateBy('')}>
        <FaClose />
      </button>
    </span>,
    <Graph key="graph" type="bar" data={data} options={options} className={style.list} />,
  ];
}

const FundsModifier = ({
  fundsSize,
  initialSize,
  filters,
  filterBy,
  order,
  orderBy,
  reverseOrder,
  aggregate,
  aggregateBy,
  onAggregateSizeChange,
  aggregated,
}) => (
  <div className={style.list}>
    {renderCount(fundsSize, initialSize)}
    {renderFilters(filters, filterBy)}
    {renderOrder(order, orderBy, reverseOrder)}
    {renderAggregate(aggregate, aggregateBy, aggregated, onAggregateSizeChange)}
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
  aggregate: PropTypes.shape({}).isRequired,
  aggregateBy: PropTypes.func.isRequired,
  onAggregateSizeChange: PropTypes.func.isRequired,
  aggregated: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
};

export default FundsModifier;
