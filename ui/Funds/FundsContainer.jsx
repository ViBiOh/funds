import React, { Component } from 'react';
import { buildFullTextRegex, fullTextRegexFilter } from '../Search/FullTextSearch';
import FundsService from '../Service/FundsService';
import setRef from '../Tools/ref';
import Throbber from '../Throbber/Throbber';
import {
  AGGREGATE_SIZES,
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM,
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM,
  RESERVED_PARAM,
} from './FundsConstantes';
import FundsHeader from './FundsHeader';
import FundsModifiers from './FundsModifiers';
import FundsGraph from './FundsGraph';
import FundsList from './FundsList';
import style from './FundsContainer.css';

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
      } else if (FundsContainer.isUndefined(o2, orderKey)) {
        return 1 * compareMultiplier;
      } else if (o1[orderKey] < o2[orderKey]) {
        return -1 * compareMultiplier;
      } else if (o1[orderKey] > o2[orderKey]) {
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
    this.setState(
      {
        aggregat: { ...this.state.aggregat, size: value.target.value },
      },
      this.filterOrderData,
    );
  }

  fetchPerformances() {
    return FundsService.getFunds()
      .then((funds) => {
        this.setState(
          {
            funds: funds.results.filter(fund => fund.id),
            loaded: true,
          },
          this.filterOrderData,
        );

        return funds;
      })
      .catch((e) => {
        console.error('Error while fetching performance:', e);
      });
  }

  filterBy(filterName, value) {
    if (value === '') {
      this.header.resetInput();
    }

    this.setState(
      {
        filters: {
          ...this.state.filters,
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
    this.setState(
      {
        order: { ...this.state.order, descending: !this.state.order.descending },
      },
      this.filterOrderData,
    );
  }

  filterOrderData() {
    const displayed = FundsContainer.filterFunds(this.state.funds, this.state.filters);

    if (this.state.order.key) {
      FundsContainer.orderFunds(displayed, this.state.order.key, this.state.order.descending);
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
    const { key } = this.state.aggregat;
    if (!key) {
      return [];
    }

    const aggregate = {};
    const size = Math.min(displayed.length, this.state.aggregat.size);
    for (let i = 0; i < size; i += 1) {
      if (typeof aggregate[displayed[i][key]] === 'undefined') {
        aggregate[displayed[i][key]] = 0;
      }
      aggregate[displayed[i][key]] += 1;
    }

    const aggregated = Object.keys(aggregate).map(label => ({
      label,
      count: aggregate[label],
    }));
    aggregated.sort((o1, o2) => o2.count - o1.count);

    return aggregated;
  }

  updateUrl() {
    const params = Object.keys(this.state.filters)
      .filter(filter => this.state.filters[filter])
      .map(filter => `${filter}=${encodeURIComponent(this.state.filters[filter])}`);

    if (this.state.order.key) {
      params.push(`${ORDER_PARAM}=${this.state.order.key}`);

      if (!this.state.order.descending) {
        params.push(ASCENDING_ORDER_PARAM);
      }
    }

    if (this.state.aggregat.key) {
      params.push(`${AGGREGAT_PARAM}=${this.state.aggregat.key}`);
      params.push(`${AGGREGAT_SIZE_PARAM}=${this.state.aggregat.size}`);
    }

    window.history.pushState(null, null, `/${params.length > 0 ? '?' : ''}${params.join('&')}`);
  }

  render() {
    return (
      <span>
        <FundsHeader
          ref={e => setRef(this, 'header', e)}
          orderBy={this.orderBy}
          aggregateBy={this.aggregateBy}
          filterBy={this.filterBy}
        />
        {this.state.error && (
          <div>
            <h2>Erreur rencontée</h2>
            <pre>{JSON.stringify(this.state.error, null, 2)}</pre>
          </div>
        )}
        <article className={style.container}>
          <div className={style.modifiers}>
            <FundsModifiers
              fundsSize={this.state.displayed.length}
              initialSize={this.state.funds.length}
              orderBy={this.orderBy}
              order={this.state.order}
              filterBy={this.filterBy}
              filters={this.state.filters}
              reverseOrder={this.reverseOrder}
              aggregateBy={this.aggregateBy}
              aggregat={this.state.aggregat}
              onAggregateSizeChange={this.onAggregateSizeChange}
            />
            <FundsGraph aggregat={this.state.aggregat} aggregated={this.state.aggregated} />
          </div>
          {!this.state.loaded && <Throbber label="Chargement des fonds" />}
          {this.state.loaded && <FundsList funds={this.state.displayed} filterBy={this.filterBy} />}
        </article>
      </span>
    );
  }
}
