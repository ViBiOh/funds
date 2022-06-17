import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import actions from 'actions';
import { isUndefined } from 'helpers/Object';
import { getSearchParamsAsObject } from 'helpers/URL';
import setRef from 'helpers/setRef';
import Throbber from 'components/Throbber';
import {
  AGGREGATE_SIZES,
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM,
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM,
  RESERVED_PARAM,
} from 'components/Funds/Constants';
import Header from 'components/Funds/Header';
import Modifiers from 'components/Funds/Modifiers';
import Graph from 'components/Funds/Graph';
import List from 'components/Funds/List';
import style from './index.module.css';

export class Funds extends Component {
  constructor(props) {
    super(props);

    this.onAggregateSizeChange = this.onAggregateSizeChange.bind(this);
    this.filterBy = this.filterBy.bind(this);
    this.aggregateBy = this.aggregateBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);
  }

  componentDidMount() {
    const { getFunds, setFilter, setOrder, setAggregat } = this.props;

    getFunds();

    const params = getSearchParamsAsObject();

    Object.entries(params)
      .filter(([key]) => !RESERVED_PARAM.includes(key))
      .forEach(([key, value]) => setFilter(key, value));

    setOrder(
      params[ORDER_PARAM] || '',
      isUndefined(params, ASCENDING_ORDER_PARAM),
    );

    setAggregat(
      params[AGGREGAT_PARAM] || '',
      parseInt(params[AGGREGAT_SIZE_PARAM], 10) || AGGREGATE_SIZES[0],
    );
  }

  onAggregateSizeChange(value) {
    const {
      setAggregat,
      funds: { aggregat },
    } = this.props;

    setAggregat(aggregat.key, parseInt(value.target.value, 10));
  }

  filterBy(filterName, value) {
    if (value === '') {
      this.header.resetInput();
    }

    const { setFilter } = this.props;
    setFilter(filterName, value);
  }

  aggregateBy(aggregat) {
    const { setAggregat } = this.props;
    setAggregat(aggregat, AGGREGATE_SIZES[0]);
  }

  orderBy(order) {
    const { setOrder } = this.props;
    setOrder(order, true);
  }

  reverseOrder() {
    const {
      setOrder,
      funds: {
        order: { key, descending },
      },
    } = this.props;
    setOrder(key, !descending);
  }

  render() {
    const {
      funds: { all, displayed, aggregated, filters, order, aggregat },
      pending,
      error,
    } = this.props;

    let content;
    if (pending) {
      content = <Throbber label="Chargement des fonds" />;
    } else {
      content = <List funds={displayed} filterBy={this.filterBy} />;
    }

    return (
      <>
        <Header
          ref={(e) => setRef(this, 'header', e)}
          orderBy={this.orderBy}
          aggregateBy={this.aggregateBy}
          filterBy={this.filterBy}
        />

        {error && (
          <div>
            <h2>Erreur rencontée</h2>
            <pre>{JSON.stringify(error, null, 2)}</pre>
          </div>
        )}

        <article className={style.container}>
          <div className={style.modifiers}>
            <Modifiers
              fundsSize={displayed.length}
              initialSize={all.length}
              orderBy={this.orderBy}
              order={order}
              filterBy={this.filterBy}
              filters={filters}
              reverseOrder={this.reverseOrder}
              aggregateBy={this.aggregateBy}
              aggregat={aggregat}
              onAggregateSizeChange={this.onAggregateSizeChange}
            />
            <Graph aggregat={aggregat} aggregated={aggregated} />
          </div>

          {content}
        </article>
      </>
    );
  }
}

Funds.propTypes = {
  error: PropTypes.string,
  getFunds: PropTypes.func.isRequired,
  pending: PropTypes.bool,
  setAggregat: PropTypes.func.isRequired,
  setFilter: PropTypes.func.isRequired,
  setOrder: PropTypes.func.isRequired,
  funds: PropTypes.shape({
    all: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
    displayed: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
    aggregated: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
    aggregat: PropTypes.shape({
      key: PropTypes.string.isRequired,
      size: PropTypes.number.isRequired,
    }).isRequired,
    order: PropTypes.shape({
      key: PropTypes.string.isRequired,
      descending: PropTypes.bool.isRequired,
    }).isRequired,
    filters: PropTypes.shape({}).isRequired,
  }).isRequired,
};

Funds.defaultProps = {
  error: null,
  pending: false,
};

/**
 * Select props from Redux state.
 * @param {Object} state Current state
 */
function mapStateToProps(state) {
  return {
    pending: state.pending[actions.GET_FUNDS],
    funds: state.funds,
  };
}

/**
 * Provide actions to dispatch.
 * @type {Object}
 */
const mapDispatchToProps = {
  getFunds: actions.getFunds,
  setFilter: actions.setFilter,
  setOrder: actions.setOrder,
  setAggregat: actions.setAggregat,
};

/**
 * Funds connected.
 */
export default connect(mapStateToProps, mapDispatchToProps)(Funds);
