import React from 'react';
import { shallow } from 'enzyme';
import Modifiers from './index';

function defaultProps() {
  return {
    fundsSize: 0,
    initialSize: 0,
    filters: {},
    filterBy: () => null,
    order: {},
    orderBy: () => null,
    reverseOrder: () => null,
    aggregat: {},
    aggregateBy: () => null,
    onAggregateSizeChange: () => null,
  };
}

it('should always render as a div', () => {
  const props = defaultProps();
  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.type()).toEqual('div');
});

it('should not render count if size match', () => {
  const props = defaultProps();
  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-count]').length).toEqual(0);
});

it('should not render filters if empty', () => {
  const props = defaultProps();
  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-filters]').length).toEqual(0);
});

it('should not render order if not specified', () => {
  const props = defaultProps();
  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-order]').length).toEqual(0);
});

it('should render count if size differ', () => {
  const props = defaultProps();
  props.fundsSize = 4000;
  props.initialSize = 8000;

  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-count]').length).toEqual(1);
  expect(wrapper.find('span[data-funds-count]').text()).toEqual('4000 / 8000');
});

it('should render each filter separately', () => {
  const props = defaultProps();
  props.filters = {
    isin: '123456',
    category: 'Test',
    score: 0,
  };

  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-filter]').length).toEqual(2);
});

it('should call given callback for removing filter', () => {
  const props = defaultProps();
  props.filters = {
    isin: '123456',
  };
  props.filterBy = jest.fn();

  const wrapper = shallow(<Modifiers {...props} />);
  wrapper
    .find('span[data-funds-filter] Button[data-funds-filter-clear]')
    .simulate('click');

  expect(props.filterBy).toHaveBeenCalledWith('isin', '');
});

it('should render order if specified', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
  };

  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-order]').length).toEqual(1);
  expect(wrapper.find('span[data-funds-order] FaSortAmountUp').length).toEqual(
    1,
  );
});

it('should render order descending if specified', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
    descending: true,
  };

  const wrapper = shallow(<Modifiers {...props} />);
  expect(
    wrapper.find('span[data-funds-order] FaSortAmountDown').length,
  ).toEqual(1);
});

it('should call given callback for removing order', () => {
  const props = defaultProps();
  props.order = {
    key: 'score',
  };
  props.orderBy = jest.fn();

  const wrapper = shallow(<Modifiers {...props} />);
  wrapper
    .find('span[data-funds-order] Button[data-funds-order-clear]')
    .simulate('click');

  expect(props.orderBy).toHaveBeenCalledWith('');
});

it('should render aggregat if specified', () => {
  const props = defaultProps();
  props.aggregat = {
    key: 'score',
  };

  const wrapper = shallow(<Modifiers {...props} />);
  expect(wrapper.find('span[data-funds-aggregat]').length).toEqual(1);
});

it('should call given callback for removing order', () => {
  const props = defaultProps();
  props.aggregat = {
    key: 'score',
  };
  props.aggregateBy = jest.fn();

  const wrapper = shallow(<Modifiers {...props} />);
  wrapper
    .find('span[data-funds-aggregat] Button[data-funds-aggregat-clear]')
    .simulate('click');

  expect(props.aggregateBy).toHaveBeenCalledWith('');
});
