import React, { Component } from 'react';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import { buildFullTextRegex, fullTextRegexFilter } from '../Search/FullTextSearch';
import Throbber from '../Throbber/Throbber';
import FundsService, { FETCH_SIZE } from './FundsService';
import FundRow from './FundRow';
import style from './Funds.css';

const COLUMNS = {
  isin: {
    label: 'ISIN',
    sortable: false,
    filterable: true,
  },
  label: {
    label: 'Libellé',
    sortable: true,
    filterable: true,
  },
  category: {
    label: 'Catégorie',
    sortable: true,
    filterable: true,
  },
  rating: {
    label: 'Note',
    sortable: true,
    filterable: true,
  },
  '1m': {
    label: '1 mois',
    sortable: true,
    filterable: false,
  },
  '3m': {
    label: '3 mois',
    sortable: true,
    filterable: false,
  },
  '6m': {
    label: '6 mois',
    sortable: true,
    filterable: false,
  },
  '1y': {
    label: '1 an',
    sortable: true,
    filterable: false,
  },
  v3y: {
    label: 'Volatilité',
    sortable: true,
    filterable: false,
  },
  score: {
    label: 'Score',
    sortable: true,
    filterable: false,
  },
};

export default class Funds extends Component {
  constructor(props) {
    super(props);

    const params = {};
    window.location.search.replace(/([^?&=]+)(?:=([^?&=]*))?/g, (match, key, value) => {
      params[key] = typeof value === 'undefined' ? true : value;
    });

    const filters = Object.assign({}, params);
    delete filters.o;
    delete filters.ao;

    this.state = {
      loaded: false,
      ids: [],
      funds: [],
      displayed: [],
      toggleDisplayed: '',
      selectedFilter: 'label',
      order: {
        key: params.o || '',
        descending: typeof params.ao === 'undefined',
      },
      filters,
    };

    this.fetchIdList = this.fetchIdList.bind(this);
    this.fetchAllPerformances = this.fetchAllPerformances.bind(this);
    this.fetchPerformances = this.fetchPerformances.bind(this);
    this.fetchPerformance = this.fetchPerformance.bind(this);

    this.onFilterChange = this.onFilterChange.bind(this);

    this.filterBy = this.filterBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);

    this.updateDataPresentation = this.updateDataPresentation.bind(this);
    this.pushHistory = this.pushHistory.bind(this);

    this.renderError = this.renderError.bind(this);
    this.renderOrderIcon = this.renderOrderIcon.bind(this);
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

  onFilterChange(e) {
    this.setState({ selectedFilter: e.target.value, toggleDisplayed: '' });
  }

  get orderDisplayed() {
    return this.state.toggleDisplayed === 'order';
  }

  set orderDisplayed(display) {
    this.setState({ toggleDisplayed: display ? 'order' : '' });
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
        }, this.updateDataPresentation);

        return funds;
      });
  }

  fetchPerformance(id) {
    return FundsService.getFund(id)
      .then((fund) => {
        this.setState({
          funds: [...this.state.funds, fund],
        }, this.updateDataPresentation);

        return fund;
      });
  }

  filterBy(filterName, value) {
    const filter = {};
    filter[filterName] = value;

    this.setState({
      filters: Object.assign(this.state.filters, filter),
    }, this.updateDataPresentation);
  }

  orderBy(order) {
    this.setState({
      order: { key: order, descending: true },
    }, this.updateDataPresentation);

    this.orderDisplayed = false;
  }

  reverseOrder() {
    this.setState({
      order: Object.assign(this.state.order, { descending: !this.state.order.descending }),
    }, this.updateDataPresentation);
  }

  updateDataPresentation() {
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
      }, this.pushHistory);
    }, 400);
  }

  pushHistory() {
    const params = Object.keys(this.state.filters)
      .filter(filter => this.state.filters[filter])
      .map(filter => `${filter}=${encodeURIComponent(this.state.filters[filter])}`);

    if (this.state.order.key) {
      params.push(`o=${this.state.order.key}`);

      if (!this.state.order.descending) {
        params.push('ao');
      }
    }

    window.history.pushState(null, null, `/?${params.join('&')}`);
  }

  renderError() {
    return (
      <div>
        <h2>Erreur rencontée</h2>
        <pre>{JSON.stringify(this.state.error, null, 2)}</pre>
      </div>
    );
  }

  renderOrderIcon() {
    const orderColumns = Object.keys(COLUMNS)
      .filter(column => COLUMNS[column].sortable)
      .map(key => (
        <li key={key}>
          <button onClick={() => this.orderBy(key)}>{COLUMNS[key].label}</button>
        </li>
      ));

    return (
      <span className={style.icon}>
        <FaSortAmountDesc
          className={this.orderDisplayed ? style.active : ''}
          onClick={() => (this.orderDisplayed = !this.orderDisplayed)}
        />
        <ol className={this.orderDisplayed ? style.displayed : style.hidden}>
          {orderColumns}
        </ol>
      </span>
    );
  }

  renderFilterIcon() {
    const filterColumns = Object.keys(COLUMNS)
      .filter(column => COLUMNS[column].filterable)
      .map(key => (
        <li key={key}>
          <button onClick={this.onFilterChange} value={key}>{COLUMNS[key].label}</button>
        </li>
      ));

    return (
      <span className={style.icon}>
        <FaFilter
          className={this.filterDisplayed ? style.active : ''}
          onClick={() => (this.filterDisplayed = !this.filterDisplayed)}
        />
        <ol className={this.filterDisplayed ? style.displayed : style.hidden}>
          {filterColumns}
        </ol>
      </span>
    );
  }

  renderHeader() {
    const count =
      `${this.state.displayed.length} / ${this.state.funds.length} / ${this.state.ids.length}`;

    return (
      <header className={style.header}>
        <h1 title={count}>Funds</h1>
        {this.renderOrderIcon()}
        {this.renderFilterIcon()}
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

  renderDataModifier() {
    return (
      <div className={style.list}>
        {this.renderFilter()}
        {this.renderOrder()}
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

Funds.propTypes = {
  location: React.PropTypes.shape({
    query: React.PropTypes.shape({
      o: React.PropTypes.string,
      ao: React.PropTypes.string,
    }),
  }).isRequired,
};
