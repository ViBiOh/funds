import React from 'react';
import PropTypes from 'prop-types';
import Throbber from '../Throbber/Throbber';
import FundsHeader from './FundsHeader';
import FundsModifier from './FundsModifier';
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
    <article>
      <FundsModifier
        fundsSize={funds.length}
        initialSize={initialSize}
        orderBy={orderBy}
        order={order}
        filterBy={filterBy}
        filters={filters}
        reverseOrder={reverseOrder}
        aggregateBy={aggregateBy}
        aggregate={aggregate}
        onAggregateSizeChange={onAggregateSizeChange}
        aggregated={aggregated}
      />
      {!loaded && <Throbber label="Chargement des fonds" />}
      {loaded && <FundsList funds={funds} filterBy={filterBy} />}
    </article>
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
