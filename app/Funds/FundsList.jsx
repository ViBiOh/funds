import React from 'react';
import PropTypes from 'prop-types';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaPieChart from 'react-icons/lib/fa/pie-chart';
import Throbber from '../Throbber/Throbber';
import { COLUMNS, CHART_COLORS, COLUMNS_HEADER, AGGREGATE_SIZES } from './FundsConstantes';
import FundRow from './FundRow';
import Graph from './Graph';
import style from './Funds.css';

function renderCount(funds, initialSize) {
  if (funds.length === initialSize) {
    return null;
  }

  return (
    <span key="count" className={style.modifier}>
      {funds.length} / {initialSize}
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
  aggregated.forEach((entry) => {
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

const FundsList = ({
  loaded,
  funds,
  aggregated,
  initialSize,
  filters,
  filterBy,
  order,
  orderBy,
  reverseOrder,
  aggregate,
  aggregateBy,
  onAggregateSizeChange,
}) => (
  <article>
    <div className={style.list}>
      {renderCount(funds, initialSize)}
      {renderFilters(filters, filterBy)}
      {renderOrder(order, orderBy, reverseOrder)}
      {renderAggregate(aggregate, aggregateBy, aggregated, onAggregateSizeChange)}
    </div>
    {!loaded && <Throbber label="Chargement des fonds" />}
    {loaded &&
      <div key="list" className={style.list}>
        <FundRow key={'header'} fund={COLUMNS_HEADER} />
        {funds.map(fund => <FundRow key={fund.id} fund={fund} filterBy={filterBy} />)}
      </div>}
  </article>
);

FundsList.displayName = 'FundsList';

FundsList.propTypes = {
  loaded: PropTypes.bool.isRequired,
  funds: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  aggregated: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  initialSize: PropTypes.number.isRequired,
  filters: PropTypes.shape({}).isRequired,
  filterBy: PropTypes.func.isRequired,
  order: PropTypes.shape({}).isRequired,
  orderBy: PropTypes.func.isRequired,
  reverseOrder: PropTypes.func.isRequired,
  aggregate: PropTypes.shape({}).isRequired,
  aggregateBy: PropTypes.func.isRequired,
  onAggregateSizeChange: PropTypes.func.isRequired,
};

export default FundsList;
