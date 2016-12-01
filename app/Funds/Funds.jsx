import React, { Component } from 'react';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaCalculator from 'react-icons/lib/fa/calculator';
import { buildFullTextRegex, fullTextRegexFilter } from '../Search/FullTextSearch';
import Throbber from '../Throbber/Throbber';
import Graph from './Graph';
import FundsService, { FETCH_SIZE } from './FundsService';
import FundRow from './FundRow';
import HeaderIcon from './HeaderIcon';
import style from './Funds.css';

const COLUMNS = {
  isin: {
    label: 'ISIN',
    sortable: true,
    summable: false,
    filterable: true,
  },
  label: {
    label: 'Libellé',
    sortable: true,
    summable: false,
    filterable: true,
  },
  category: {
    label: 'Catégorie',
    sortable: true,
    summable: true,
    filterable: true,
  },
  rating: {
    label: 'Note',
    sortable: true,
    summable: true,
    filterable: true,
  },
  '1m': {
    label: '1 mois',
    sortable: true,
    summable: false,
    filterable: false,
  },
  '3m': {
    label: '3 mois',
    sortable: true,
    summable: false,
    filterable: false,
  },
  '6m': {
    label: '6 mois',
    sortable: true,
    summable: false,
    filterable: false,
  },
  '1y': {
    label: '1 an',
    sortable: true,
    summable: false,
    filterable: false,
  },
  v3y: {
    label: 'Volatilité',
    sortable: true,
    summable: false,
    filterable: false,
  },
  score: {
    label: 'Score',
    sortable: true,
    summable: false,
    filterable: false,
  },
};

const CHART_COLORS = [
  '#1f77b4', '#e377c2', '#ff7f0e', '#2ca02c', '#bcbd22', '#d62728',
  '#17becf', '#9467bd', '#7f7f7f', '#8c564b', '#3366cc',
];

const SUM_PARAM = 's';
const ORDER_PARAM = 'o';
const ASCENDING_ORDER_PARAM = 'ao';
const RESERVED_PARAM = [SUM_PARAM, ORDER_PARAM, ASCENDING_ORDER_PARAM];

export default class Funds extends Component {
  constructor(props) {
    super(props);

    const params = {};
    window.location.search.replace(/([^?&=]+)(?:=([^?&=]*))?/g, (match, key, value) => {
      params[key] = typeof value === 'undefined' ? true : decodeURIComponent(value);
    });

    const filters = Object.assign({}, params);
    RESERVED_PARAM.forEach(param => delete filters[param]);

    this.state = {
      loaded: false,
      ids: [],
      funds: [],
      displayed: [],
      aggregated: [],
      toggleDisplayed: '',
      selectedFilter: 'label',
      sum: {
        key: params[SUM_PARAM] || '',
        size: 25,
      },
      order: {
        key: params[ORDER_PARAM] || '',
        descending: typeof params[ASCENDING_ORDER_PARAM] === 'undefined',
      },
      filters,
    };

    this.fetchIdList = this.fetchIdList.bind(this);
    this.fetchAllPerformances = this.fetchAllPerformances.bind(this);
    this.fetchPerformances = this.fetchPerformances.bind(this);
    this.fetchPerformance = this.fetchPerformance.bind(this);

    this.onFilterChange = this.onFilterChange.bind(this);

    this.filterBy = this.filterBy.bind(this);
    this.aggregateBy = this.aggregateBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);

    this.updateData = this.updateData.bind(this);
    this.aggregateData = this.aggregateData.bind(this);
    this.pushHistory = this.pushHistory.bind(this);

    this.renderError = this.renderError.bind(this);
    this.renderFilterIcon = this.renderFilterIcon.bind(this);
    this.renderHeader = this.renderHeader.bind(this);

    this.renderFilter = this.renderFilter.bind(this);
    this.renderOrder = this.renderOrder.bind(this);
    this.renderDataModifier = this.renderDataModifier.bind(this);

    this.renderList = this.renderList.bind(this);
    this.renderRow = this.renderRow.bind(this);
  }

  componentDidMount() {
    this.fetchIdList()
      .then(this.fetchAllPerformances)
      .catch(error => this.setState({ error }));
  }

  onFilterChange(selectedFilter) {
    this.setState({ selectedFilter, toggleDisplayed: '' });
  }

  get orderDisplayed() {
    return this.state.toggleDisplayed === 'order';
  }

  set orderDisplayed(display) {
    this.setState({ toggleDisplayed: display ? 'order' : '' });
  }

  get sigmaDisplayed() {
    return this.state.toggleDisplayed === 'sigma';
  }

  set sigmaDisplayed(display) {
    this.setState({ toggleDisplayed: display ? 'sigma' : '' });
  }

  get filterDisplayed() {
    return this.state.toggleDisplayed === 'filter';
  }

  set filterDisplayed(display) {
    this.setState({ toggleDisplayed: display ? 'filter' : '' });
  }

  fetchIdList() {
    return FundsService.getIdList()
      .then((ids) => {
        this.setState({ ids });
        return ids;
      });
  }

  fetchAllPerformances() {
    const fetches = [];
    for (let i = 0, size = this.state.ids.length; i < size; i += FETCH_SIZE) {
      fetches.push(this.fetchPerformances(this.state.ids.slice(i, i + FETCH_SIZE)));
    }

    Promise.all(fetches).then(() => this.setState({ loaded: true }));
  }

  fetchPerformances(ids) {
    return FundsService.getFunds(ids)
      .then((funds) => {
        const results = funds.results.filter(fund => fund.id);
        this.setState({
          funds: [...this.state.funds, ...results],
        }, this.updateData);

        return funds;
      });
  }

  fetchPerformance(id) {
    return FundsService.getFund(id)
      .then((fund) => {
        this.setState({
          funds: [...this.state.funds, fund],
        }, this.updateData);

        return fund;
      });
  }

  filterBy(filterName, value) {
    const filter = {};
    filter[filterName] = value;

    this.setState({
      filters: Object.assign(this.state.filters, filter),
    }, this.updateData);
  }

  aggregateBy(sum) {
    this.setState({
      sum: { key: sum, size: 25 },
    }, this.updateData);

    this.sigmaDisplayed = false;
  }

  orderBy(order) {
    this.setState({
      order: { key: order, descending: true },
    }, this.updateData);

    this.orderDisplayed = false;
  }

  reverseOrder() {
    this.setState({
      order: Object.assign(this.state.order, { descending: !this.state.order.descending }),
    }, this.updateData);
  }

  updateData() {
    clearTimeout(this.timeout);
    this.timeout = setTimeout(() => {
      const displayed = Object.keys(this.state.filters).reduce((previous, filter) => {
        const regex = buildFullTextRegex(this.state.filters[filter]);
        return previous.filter(fund => fullTextRegexFilter(fund[filter], regex));
      }, this.state.funds.slice());

      if (this.state.order.key) {
        const orderKey = this.state.order.key;
        const compareMultiplier = this.state.order.descending ? -1 : 1;

        displayed.sort((o1, o2) => {
          if (!o1 || typeof o1[orderKey] === 'undefined') {
            return -1 * compareMultiplier;
          } else if (!o2 || typeof o2[orderKey] === 'undefined') {
            return 1 * compareMultiplier;
          } else if (o1[orderKey] < o2[orderKey]) {
            return -1 * compareMultiplier;
          } else if (o1[orderKey] > o2[orderKey]) {
            return 1 * compareMultiplier;
          }
          return 0;
        });
      }

      this.setState({
        displayed,
        aggregated: this.aggregateData(displayed),
      }, this.pushHistory);
    }, 400);
  }

  aggregateData(displayed) {
    const { key } = this.state.sum;
    if (!key) {
      return [];
    }

    const aggregate = {};
    const size = Math.min(displayed.length, this.state.sum.size);
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

  pushHistory() {
    const params = Object.keys(this.state.filters)
      .filter(filter => this.state.filters[filter])
      .map(filter => `${filter}=${encodeURIComponent(this.state.filters[filter])}`);

    if (this.state.order.key) {
      params.push(`${ORDER_PARAM}=${this.state.order.key}`);

      if (!this.state.order.descending) {
        params.push(ASCENDING_ORDER_PARAM);
      }
    }

    if (this.state.sum.key) {
      params.push(`${SUM_PARAM}=${this.state.sum.key}`);
    }

    window.history.pushState(null, null, `/${params.length > 0 ? '?' : ''}${params.join('&')}`);
  }

  renderError() {
    return (
      <div>
        <h2>Erreur rencontée</h2>
        <pre>{JSON.stringify(this.state.error, null, 2)}</pre>
      </div>
    );
  }

  renderHeader() {
    return (
      <header className={style.header}>
        <h1>Funds</h1>
        <HeaderIcon
          columns={COLUMNS}
          filter="sortable"
          onClick={this.orderBy}
          icon={
            <FaSortAmountDesc
              onClick={() => (this.orderDisplayed = !this.orderDisplayed)}
            />
          }
          displayed={this.orderDisplayed}
        />
        <HeaderIcon
          columns={COLUMNS}
          filter="summable"
          onClick={this.aggregateBy}
          icon={
            <FaCalculator
              onClick={() => (this.sigmaDisplayed = !this.sigmaDisplayed)}
            />
          }
          displayed={this.sigmaDisplayed}
        />
        <HeaderIcon
          columns={COLUMNS}
          filter="filterable"
          onClick={this.onFilterChange}
          icon={
            <FaFilter
              onClick={() => (this.filterDisplayed = !this.filterDisplayed)}
            />
          }
          displayed={this.filterDisplayed}
        />
        <input
          type="text"
          placeholder={`Fitre sur ${COLUMNS[this.state.selectedFilter].label}`}
          onChange={e => this.filterBy(this.state.selectedFilter, e.target.value)}
        />
        {!this.state.loaded && <Throbber />}
      </header>
    );
  }

  renderFilter() {
    return Object.keys(this.state.filters)
      .filter(filter => this.state.filters[filter])
      .map(filter => (
        <span key={filter} className={style.dataModifier}>
          <span className={style.icon}>
            <FaFilter />
          </span>
          <span><em> {COLUMNS[filter].label}</em> &#x2243; </span>
          {this.state.filters[filter]}
          <button onClick={() => this.filterBy(filter, '')} className={style.icon}>
            <FaClose />
          </button>
        </span>
      ));
  }

  renderOrder() {
    return this.state.order.key && (
      <span className={style.dataModifier}>
        <button onClick={this.reverseOrder} className={style.icon}>
          {this.state.order.descending ? <FaSortAmountDesc /> : <FaSortAmountAsc />}
        </button>
        &nbsp;{COLUMNS[this.state.order.key].label}
        <button onClick={() => this.orderBy('')} className={style.icon}>
          <FaClose />
        </button>
      </span>
    );
  }

  renderSigma() {
    if (!this.state.sum.key) {
      return null;
    }

    const { aggregated } = this.state;
    const label = COLUMNS[this.state.sum.key].label;

    const options = {
      legend: false,
      scales: {
        xAxes: [{
          display: false,
        }],
        yAxes: [{
          display: false,
          ticks: {
            beginAtZero: true,
          },
        }],
      },
    };

    const data = {
      labels: [],
      datasets: [{
        label: `\u03A3 ${this.state.sum.size} | ${label}`,
        data: [],
        backgroundColor: [],
      }],
    };

    let i = 0;
    aggregated.forEach((entry) => {
      data.labels.push(entry.label);
      data.datasets[0].data.push(entry.count);
      data.datasets[0].backgroundColor.push(CHART_COLORS[i]);

      i = (i + 1) % CHART_COLORS.length;
    });

    return [
      <span key="label" className={style.dataModifier}>
        &#x3A3; {this.state.sum.size} | {label}
        <button onClick={() => this.aggregateBy('')} className={style.icon}>
          <FaClose />
        </button>
      </span>,
      <Graph
        key="graph"
        type="bar"
        data={data}
        options={options}
        className={style.list}
      />,
    ];
  }

  renderDataModifier() {
    return (
      <div className={style.list}>
        {this.renderFilter()}
        {this.renderOrder()}
        {this.renderSigma()}
      </div>
    );
  }

  renderRow() {
    return this.state.displayed.map(fund => (
      <FundRow key={fund.id} fund={fund} filterBy={this.filterBy} />
    ));
  }

  renderList() {
    const header = Object.keys(COLUMNS).reduce((previous, current) => {
      previous[current] = COLUMNS[current].label; // eslint-disable-line no-param-reassign
      return previous;
    }, {});

    return (
      <div key="list" className={style.list}>
        <FundRow key={'header'} fund={header} />
        {this.renderRow()}
      </div>
    );
  }

  renderContent() {
    return (
      <article>
        {this.renderDataModifier()}
        {this.renderList()}
      </article>
    );
  }

  render() {
    return (
      <span>
        {this.renderHeader()}
        {
          this.state.error && this.renderError()
        }
        {this.renderContent()}
      </span>
    );
  }
}
