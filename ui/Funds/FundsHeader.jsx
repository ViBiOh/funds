import React, { Component } from 'react';
import PropTypes from 'prop-types';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import FaPieChart from 'react-icons/lib/fa/pie-chart';
import FaFilter from 'react-icons/lib/fa/filter';
import setRef from '../Tools/ref';
import { COLUMNS } from './FundsConstantes';
import HeaderIcon from './HeaderIcon';
import style from './FundsHeader.less';

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
    this.onTextChangeDebounce = this.onTextChangeDebounce.bind(this);
    this.toggleDisplay = this.toggleDisplay.bind(this);
    this.resetInput = this.resetInput.bind(this);
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
    this.resetInput();
    this.setState({ selectedFilter, toggleDisplayed: '' });
  }

  onTextChangeDebounce(e) {
    clearTimeout(this.timeout);
    ((text) => {
      this.timeout = setTimeout(() => this.props.filterBy(this.state.selectedFilter, text), 300);
      return undefined;
    })(e.target.value);
  }

  get orderDisplayed() {
    return this.state.toggleDisplayed === 'order';
  }

  get sigmaDisplayed() {
    return this.state.toggleDisplayed === 'sigma';
  }

  get filterDisplayed() {
    return this.state.toggleDisplayed === 'filter';
  }

  toggleDisplay(icon, display) {
    this.setState({ toggleDisplayed: display ? icon : '' });
  }

  resetInput() {
    this.filter.value = '';
  }

  render() {
    return (
      <header className={style.header}>
        <h1>Funds</h1>
        <HeaderIcon
          filter="sortable"
          onClick={this.onOrderBy}
          icon={
            <FaSortAmountDesc onClick={() => this.toggleDisplay('order', !this.orderDisplayed)} />
          }
          displayed={this.orderDisplayed}
        />
        <HeaderIcon
          filter="summable"
          onClick={this.onAggregateBy}
          icon={<FaPieChart onClick={() => this.toggleDisplay('sigma', !this.sigmaDisplayed)} />}
          displayed={this.sigmaDisplayed}
        />
        <HeaderIcon
          filter="filterable"
          onClick={this.onFilterChange}
          icon={<FaFilter onClick={() => this.toggleDisplay('filter', !this.filterDisplayed)} />}
          displayed={this.filterDisplayed}
        />
        <input
          type="text"
          ref={e => setRef(this, 'filter', e)}
          placeholder={`Fitre sur ${COLUMNS[this.state.selectedFilter].label}`}
          onChange={this.onTextChangeDebounce}
        />
      </header>
    );
  }
}

FundsHeader.propTypes = {
  aggregateBy: PropTypes.func.isRequired,
  filterBy: PropTypes.func.isRequired,
  orderBy: PropTypes.func.isRequired,
};
