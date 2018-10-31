import React from 'react';
import { shallow } from 'enzyme';
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

test('should always render as a span', () => {
  const wrapper = shallow(<FundRow {...defaultProps} />);

  wrapper
    .find('Button')
    .at(0)
    .simulate('click');

  expect(wrapper.type()).toEqual('span');
});

test('should call given filterBy func on category click', () => {
  const filterBy = jest.fn();
  const wrapper = shallow(<FundRow {...defaultProps} filterBy={filterBy} />);

  wrapper
    .find('Button')
    .at(0)
    .simulate('click');

  expect(filterBy).toHaveBeenCalledWith('category', 'Test');
});

test('should call given filterBy func on rating click', () => {
  const filterBy = jest.fn();
  const wrapper = shallow(<FundRow {...defaultProps} filterBy={filterBy} />);

  wrapper
    .find('Button')
    .at(1)
    .simulate('click');

  expect(filterBy).toHaveBeenCalledWith('rating', 4);
});
