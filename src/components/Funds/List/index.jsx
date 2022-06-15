import React from 'react';
import PropTypes from 'prop-types';
import { COLUMNS_HEADER } from 'components/Funds/Constants';
import Row from 'components/Funds/Row';

/**
 * List functional component.
 * @param {Array}    options.funds     List of Funds to render
 * @param {Function} options.filterBy  Filter function
 */
export default function List({ funds, filterBy }) {
  return (
    <div key="list">
      <Row key="header" fund={COLUMNS_HEADER} dataTestId="row-header" />
      {funds.map((fund) => (
        <Row
          key={fund.id}
          fund={fund}
          filterBy={filterBy}
          dataTestId="fund-row"
        />
      ))}
    </div>
  );
}

List.displayName = 'List';

List.propTypes = {
  funds: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
    }),
  ).isRequired,
  filterBy: PropTypes.func.isRequired,
};
