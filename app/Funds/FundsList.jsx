import React from 'react';
import PropTypes from 'prop-types';
import { COLUMNS_HEADER } from './FundsConstantes';
import FundRow from './FundRow';
import style from './FundsList.css';

const FundsList = ({ funds, filterBy }) => (
  <div key="list" className={style.list}>
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
