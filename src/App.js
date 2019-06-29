import React, { Component } from 'react';
import { buildFullTextRegex, fullTextRegexFilter } from 'helpers/Search';
import FundsService from 'services/Funds';
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
import FundsHeader from 'containers/Funds/FundsHeader';
import FundsModifiers from 'containers/Funds/FundsModifiers';
import FundsGraph from 'containers/Funds/FundsGraph';
import List from 'components/Funds/List';
import style from './App.module.css';

export default class FundsContainer extends Component {
  static isUndefined(o, orderKey) {
    return !o || typeof o[orderKey] === 'undefined';
  }

  static filterFunds(funds, filters) {
    return Object.keys(filters).reduce((previous, filter) => {
      const regex = buildFullTextRegex(filters[filter]);
      return previous.filter(fund => fullTextRegexFilter(fund[filter], regex));
    }, funds.slice());
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

    const params = {};
    window.location.search.replace(/([^?&=]+)(?:=([^?&=]*))?/g, (match, key, value) => {
      params[key] = typeof value === 'undefined' ? true : decodeURIComponent(value);
    });

    const filters = { ...params };
    RESERVED_PARAM.forEach(param => delete filters[param]);

    this.state = {
      loaded: false,
      funds: [],
      displayed: [],
      aggregated: [],
      aggregat: {
        key: params[AGGREGAT_PARAM] || '',
        size: params[AGGREGAT_SIZE_PARAM] || AGGREGATE_SIZES[0],
      },
      order: {
        key: params[ORDER_PARAM] || '',
        descending: typeof params[ASCENDING_ORDER_PARAM] === 'undefined',
      },
      filters,
    };

    this.fetchPerformances = this.fetchPerformances.bind(this);
    this.onAggregateSizeChange = this.onAggregateSizeChange.bind(this);

    this.filterBy = this.filterBy.bind(this);
    this.aggregateBy = this.aggregateBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);

    this.filterOrderData = this.filterOrderData.bind(this);
    this.aggregateData = this.aggregateData.bind(this);
    this.updateUrl = this.updateUrl.bind(this);
  }

  componentDidMount() {
    this.fetchPerformances();
  }

  onAggregateSizeChange(value) {
    const { aggregat } = this.state;

    this.setState(
      {
        aggregat: { ...aggregat, size: value.target.value },
      },
      this.filterOrderData,
    );
  }

  fetchPerformances() {
    return FundsService.getFunds()
      .then(funds => {
        this.setState(
          {
            funds: funds.results.filter(fund => fund.id),
            loaded: true,
          },
          this.filterOrderData,
        );

        return funds;
      })
      .catch(e => {
        global.console.error('Error while fetching performance:', e);
      });
  }

  filterBy(filterName, value) {
    if (value === '') {
      this.header.resetInput();
    }

    const { filters } = this.state;

    this.setState(
      {
        filters: {
          ...filters,
          [filterName]: value,
        },
      },
      this.filterOrderData,
    );
  }

  aggregateBy(aggregat) {
    this.setState(
      {
        aggregat: { key: aggregat, size: 25 },
      },
      this.filterOrderData,
    );
  }

  orderBy(order) {
    this.setState(
      {
        order: { key: order, descending: true },
      },
      this.filterOrderData,
    );
  }

  reverseOrder() {
    const { order } = this.state;

    this.setState(
      {
        order: { ...order, descending: !order.descending },
      },
      this.filterOrderData,
    );
  }

  filterOrderData() {
    const { funds, filters, order } = this.state;

    const displayed = FundsContainer.filterFunds(funds, filters);

    if (order.key) {
      FundsContainer.orderFunds(displayed, order.key, order.descending);
    }

    this.setState(
      {
        displayed,
        aggregated: this.aggregateData(displayed),
      },
      this.updateUrl,
    );
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
    const { error, displayed, funds, order, filters, aggregat, aggregated, loaded } = this.state;

    return (
      <>
        <FundsHeader
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
            <FundsModifiers
              fundsSize={displayed.length}
              initialSize={funds.length}
              orderBy={this.orderBy}
              order={order}
              filterBy={this.filterBy}
              filters={filters}
              reverseOrder={this.reverseOrder}
              aggregateBy={this.aggregateBy}
              aggregat={aggregat}
              onAggregateSizeChange={this.onAggregateSizeChange}
            />
            <FundsGraph aggregat={aggregat} aggregated={aggregated} />
          </div>
          {!loaded && <Throbber label="Chargement des fonds" />}
          {loaded && <List funds={displayed} filterBy={this.filterBy} />}
        </article>
      </>
    );
  }
}
