import React from 'react';
import PropTypes from 'prop-types';
import { COLUMNS_HEADER } from 'containers/Funds/FundsConstantes';
import FundRow from 'containers/Funds/FundRow';

const FundsList = ({ funds, filterBy }) => (
  <div key="list">
    <FundRow key="header" fund={COLUMNS_HEADER} />
    {funds.map(fund => <FundRow key={fund.id} fund={fund} filterBy={filterBy} />)}
  </div>
);

FundsList.displayName = 'FundsList';

FundsList.propTypes = {
  funds: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  filterBy: PropTypes.func.isRequired,
};

export default FundsList;
