import React from 'react';
import { shallow } from 'enzyme';
import Row from './index';

function defaultProps() {
  return {
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
}

it('should always render as a span', () => {
  const props = defaultProps();
  const wrapper = shallow(<Row {...props} />);

  wrapper
    .find('Button')
    .at(0)
    .simulate('click');

  expect(wrapper.type()).toEqual('span');
});

it('should call given filterBy func on category click', () => {
  const props = defaultProps();
  props.filterBy = jest.fn();
  const wrapper = shallow(<Row {...props} />);

  wrapper
    .find('Button')
    .at(0)
    .simulate('click');

  expect(props.filterBy).toHaveBeenCalledWith('category', 'Test');
});

it('should call given filterBy func on rating click', () => {
  const props = defaultProps();
  props.filterBy = jest.fn();
  const wrapper = shallow(<Row {...props} />);

  wrapper
    .find('Button')
    .at(1)
    .simulate('click');

  expect(props.filterBy).toHaveBeenCalledWith('rating', 4);
});
