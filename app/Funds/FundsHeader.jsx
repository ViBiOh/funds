import React, { Component } from 'react';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaCalculator from 'react-icons/lib/fa/calculator';
import FaFilter from 'react-icons/lib/fa/filter';
import HeaderIcon from './HeaderIcon';
import style from './Funds.css';

export default class FundsHeader extends Component {
  constructor(props) {
    super(props);

    this.state = {
      toggleDisplayed: '',
      selectedFilter: 'label',
    };

    this.onOrderBy = this.onOrderBy.bind(this);
    this.onAggregateBy = this.onAggregateBy.bind(this);
    this.onFilterChange = this.onFilterChange.bind(this);
  }

  onOrderBy(...args) {
    this.props.orderBy(...args);
    this.setState({ toggleDisplayed: '' });
  }

  onAggregateBy(...args) {
    this.props.aggregateBy(...args);
    this.setState({ toggleDisplayed: '' });
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

  render() {
    return (
      <header className={style.header}>
        <h1>Funds</h1>
        <HeaderIcon
          columns={this.props.columns}
          filter="sortable"
          onClick={this.onOrderBy}
          icon={
            <FaSortAmountDesc
              onClick={() => (this.orderDisplayed = !this.orderDisplayed)}
            />
          }
          displayed={this.orderDisplayed}
        />
        <HeaderIcon
          columns={this.props.columns}
          filter="summable"
          onClick={this.onAggregateBy}
          icon={
            <FaCalculator
              onClick={() => (this.sigmaDisplayed = !this.sigmaDisplayed)}
            />
          }
          displayed={this.sigmaDisplayed}
        />
        <HeaderIcon
          columns={this.props.columns}
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
          placeholder={`Fitre sur ${this.props.columns[this.state.selectedFilter].label}`}
          onChange={e => this.props.filterBy(this.state.selectedFilter, e.target.value)}
        />
      </header>
    );
  }
}

FundsHeader.propTypes = {
  columns: React.PropTypes.shape({}).isRequired,
  orderBy: React.PropTypes.func.isRequired,
  aggregateBy: React.PropTypes.func.isRequired,
  filterBy: React.PropTypes.func.isRequired,
};
