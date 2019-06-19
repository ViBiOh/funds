import React from 'react';
import PropTypes from 'prop-types';
import { COLUMNS_HEADER } from 'containers/Funds/FundsConstantes';
import Row from 'components/Fund/Row';

const FundsList = ({ funds, filterBy }) => (
  <div key="list">
    <Row key="header" fund={COLUMNS_HEADER} />
    {funds.map(fund => (
      <Row key={fund.id} fund={fund} filterBy={filterBy} />
    ))}
  </div>
);

FundsList.displayName = 'FundsList';

FundsList.propTypes = {
  funds: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  filterBy: PropTypes.func.isRequired,
};

export default FundsList;
