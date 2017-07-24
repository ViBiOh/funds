import test from 'ava';
import React from 'react';
import { shallow } from 'enzyme';
import sinon from 'sinon';
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
  const wrapper = shallow(<FundRow {...defaultProps} />);

  wrapper.find('Button').at(0).simulate('click');

  t.is(wrapper.type(), 'span');
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
