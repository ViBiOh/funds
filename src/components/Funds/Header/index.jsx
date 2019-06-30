import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { FaSortAmountDown, FaChartPie, FaFilter } from 'react-icons/fa';
import setRef from 'helpers/ref';
import { COLUMNS } from 'components/Funds/Constants';
import HeaderIcon from 'components/Funds/HeaderIcon';
import style from './index.module.css';

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
    const { orderBy } = this.props;
    orderBy(...args);
    this.setState({ toggleDisplayed: '' });
  }

  onAggregateBy(...args) {
    const { aggregateBy } = this.props;
    aggregateBy(...args);
    this.setState({ toggleDisplayed: '' });
  }

  onFilterChange(selectedFilter) {
    this.resetInput();
    this.setState({ selectedFilter, toggleDisplayed: '' });
  }

  onTextChangeDebounce(e) {
    clearTimeout(this.timeout);
    (text => {
      const { filterBy } = this.props;
      const { selectedFilter } = this.state;

      this.timeout = setTimeout(() => filterBy(selectedFilter, text), 300);
      return undefined;
    })(e.target.value);
  }

  get orderDisplayed() {
    const { toggleDisplayed } = this.state;

    return toggleDisplayed === 'order';
  }

  get sigmaDisplayed() {
    const { toggleDisplayed } = this.state;

    return toggleDisplayed === 'sigma';
  }

  get filterDisplayed() {
    const { toggleDisplayed } = this.state;

    return toggleDisplayed === 'filter';
  }

  toggleDisplay(icon, display) {
    this.setState({ toggleDisplayed: display ? icon : '' });
  }

  resetInput() {
    this.filter.value = '';
  }

  render() {
    const { selectedFilter } = this.state;

    return (
      <header className={style.header}>
        <h1>Funds</h1>
        <HeaderIcon
          filter="sortable"
          onClick={this.onOrderBy}
          icon={
            <FaSortAmountDown onClick={() => this.toggleDisplay('order', !this.orderDisplayed)} />
          }
          displayed={this.orderDisplayed}
        />
        <HeaderIcon
          filter="summable"
          onClick={this.onAggregateBy}
          icon={<FaChartPie onClick={() => this.toggleDisplay('sigma', !this.sigmaDisplayed)} />}
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
          placeholder={`Fitre sur ${COLUMNS[selectedFilter].label}`}
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
