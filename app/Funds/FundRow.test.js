import test from 'ava';
import React from 'react';
import { shallow } from 'enzyme';
import sinon from 'sinon';
import FundsService from '../Service/FundsService';
import FundRow from './FundRow';

const defaultProps = {
  fund: {
    isin: 'ABCDEFG',
    label: 'Fund with high risk',
    category: 'Test',
    rating: 4,
    '1m': 2.4,
    '3m': 3.4,
    '6m': 4.4,
    '1y': 5.5,
    v3y: 6.6,
    score: 25,
  },
};

test('should always render as a span', (t) => {
  t.is(shallow(<FundRow {...defaultProps} />).type(), 'span');
});

test('should set href if id is provided', (t) => {
  const wrapper = shallow(<FundRow {...defaultProps} fund={{ ...defaultProps.fund, id: 8000 }} />);

  wrapper.find('Button').at(0).simulate('click');

  t.is(wrapper.findWhere(e => e.props().href === FundsService.getDataUrl(8000)).length, 1);
});

test('should display not fresh data in special color', (t) => {
  const eightHoursAgo = new Date();
  eightHoursAgo.setTime(new Date().getTime() - 28800000);
  const wrapper = shallow(
    <FundRow {...defaultProps} fund={{ ...defaultProps.fund, ts: eightHoursAgo.toISOString() }} />,
  );

  t.is(wrapper.find('a').props().className.split(' ').length, 2);
});

test('should call given filterBy func on category click', (t) => {
  const filterBy = sinon.spy();
  const wrapper = shallow(<FundRow {...defaultProps} filterBy={filterBy} />);

  wrapper.find('Button').at(0).simulate('click');

  t.truthy(filterBy.calledWith('category', 'Test'));
});

test('should call given filterBy func on rating click', (t) => {
  const filterBy = sinon.spy();
  const wrapper = shallow(<FundRow {...defaultProps} filterBy={filterBy} />);

  wrapper.find('Button').at(1).simulate('click');

  t.truthy(filterBy.calledWith('rating', 4));
});
