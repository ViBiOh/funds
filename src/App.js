import React, { Component } from 'react';
import { connect } from 'react-redux';
import actions from 'actions';
import { getSearchParamsAsObject } from 'helpers/URL';
import setRef from 'helpers/ref';
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
import style from './App.module.css';

export class FundsContainer extends Component {
  static isUndefined(o, orderKey) {
    return !o || typeof o[orderKey] === 'undefined';
  }

  static orderFunds(funds, orderKey, reverse) {
    const compareMultiplier = reverse ? -1 : 1;

    funds.sort((o1, o2) => {
      if (FundsContainer.isUndefined(o1, orderKey)) {
        return -1 * compareMultiplier;
      }
      if (FundsContainer.isUndefined(o2, orderKey)) {
        return 1 * compareMultiplier;
      }
      if (o1[orderKey] < o2[orderKey]) {
        return -1 * compareMultiplier;
      }
      if (o1[orderKey] > o2[orderKey]) {
        return 1 * compareMultiplier;
      }
      return 0;
    });
  }

  constructor(props) {
    super(props);

    const params = getSearchParamsAsObject();

    this.state = {
      aggregated: [],
      aggregat: {
        key: params[AGGREGAT_PARAM] || '',
        size: params[AGGREGAT_SIZE_PARAM] || AGGREGATE_SIZES[0],
      },
      order: {
        key: params[ORDER_PARAM] || '',
        descending: typeof params[ASCENDING_ORDER_PARAM] === 'undefined',
      },
    };

    this.onAggregateSizeChange = this.onAggregateSizeChange.bind(this);

    this.filterBy = this.filterBy.bind(this);
    this.aggregateBy = this.aggregateBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);

    this.updateUrl = this.updateUrl.bind(this);
  }

  componentDidMount() {
    this.props.getFunds();

    Object.entries(getSearchParamsAsObject())
      .filter(([key]) => !RESERVED_PARAM.includes(key))
      .forEach(([key, value]) => this.props.setFilter(key, value));
  }

  onAggregateSizeChange(value) {
    const { aggregat } = this.state;

    this.setState({
      aggregat: { ...aggregat, size: value.target.value },
    });
  }

  filterBy(filterName, value) {
    if (value === '') {
      this.header.resetInput();
    }

    this.props.setFilter(filterName, value);
  }

  aggregateBy(aggregat) {
    this.setState({
      aggregat: { key: aggregat, size: 25 },
    });
  }

  orderBy(order) {
    this.setState({
      order: { key: order, descending: true },
    });
  }

  reverseOrder() {
    const { order } = this.state;

    this.setState({
      order: { ...order, descending: !order.descending },
    });
  }

  aggregateData(displayed) {
    const { aggregat } = this.state;

    if (!aggregat.key) {
      return [];
    }

    const aggregate = {};
    const size = Math.min(displayed.length, aggregat.size);
    for (let i = 0; i < size; i += 1) {
      if (typeof aggregate[displayed[i][aggregat.key]] === 'undefined') {
        aggregate[displayed[i][aggregat.key]] = 0;
      }
      aggregate[displayed[i][aggregat.key]] += 1;
    }

    const aggregated = Object.keys(aggregate).map(label => ({
      label,
      count: aggregate[label],
    }));
    aggregated.sort((o1, o2) => o2.count - o1.count);

    return aggregated;
  }

  updateUrl() {
    const { filters, order, aggregat } = this.state;

    const params = Object.keys(filters)
      .filter(filter => filters[filter])
      .map(filter => `${filter}=${encodeURIComponent(filters[filter])}`);

    if (order.key) {
      params.push(`${ORDER_PARAM}=${order.key}`);

      if (!order.descending) {
        params.push(ASCENDING_ORDER_PARAM);
      }
    }

    if (aggregat.key) {
      params.push(`${AGGREGAT_PARAM}=${aggregat.key}`);
      params.push(`${AGGREGAT_SIZE_PARAM}=${aggregat.size}`);
    }

    window.history.pushState(null, null, `/${params.length > 0 ? '?' : ''}${params.join('&')}`);
  }

  render() {
    const {
      funds: { filters, all, displayed },
      pending,
    } = this.props;
    const { error, order, aggregat, aggregated } = this.state;

    return (
      <>
        <Header
          ref={e => setRef(this, 'header', e)}
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
          {pending && <Throbber label="Chargement des fonds" />}
          {!pending && <List funds={displayed} filterBy={this.filterBy} />}
        </article>
      </>
    );
  }
}

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
};

/**
 * FundsContainer connected.
 */
export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(FundsContainer);
