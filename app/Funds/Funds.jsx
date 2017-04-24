import React from 'react';
import PropTypes from 'prop-types';
import FundsHeader from './FundsHeader';
import FundsList from './FundsList';

const Funds = ({
  orderBy,
  order,
  reverseOrder,
  aggregateBy,
  aggregate,
  onAggregateSizeChange,
  filterBy,
  filters,
  error,
  loaded,
  funds,
  aggregated,
  initialSize,
}) => (
  <span>
    <FundsHeader orderBy={orderBy} aggregateBy={aggregateBy} filterBy={filterBy} />
    {error &&
      <div>
        <h2>Erreur rencontée</h2>
        <pre>{JSON.stringify(error, null, 2)}</pre>
      </div>}
    <FundsList
      loaded={loaded}
      funds={funds}
      orderBy={orderBy}
      order={order}
      reverseOrder={reverseOrder}
      aggregateBy={aggregateBy}
      aggregate={aggregate}
      onAggregateSizeChange={onAggregateSizeChange}
      filterBy={filterBy}
      filters={filters}
      aggregated={aggregated}
      initialSize={initialSize}
    />
  </span>
);

Funds.displayName = 'Funds';

Funds.propTypes = {
  orderBy: PropTypes.func.isRequired,
  order: PropTypes.shape({}).isRequired,
  reverseOrder: PropTypes.func.isRequired,
  aggregateBy: PropTypes.func.isRequired,
  aggregate: PropTypes.shape({}).isRequired,
  onAggregateSizeChange: PropTypes.func.isRequired,
  filterBy: PropTypes.func.isRequired,
  filters: PropTypes.shape({}).isRequired,
  error: PropTypes.shape({}),
  loaded: PropTypes.bool.isRequired,
  funds: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  aggregated: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  initialSize: PropTypes.number.isRequired,
};

Funds.defaultProps = {
  error: undefined,
};

export default Funds;
