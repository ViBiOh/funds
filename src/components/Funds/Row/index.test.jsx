import React from 'react';
import { shallow } from 'enzyme';
import Row from './index';

function defaultProps() {
  return {
    fund: { category: 'Test', rating: 4, '1m': 0.0, '3m': 0.0, '6m': 0.0, '1y': 0.0, v3y: 0.0, score: 0.0 },
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
