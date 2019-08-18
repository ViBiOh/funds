import React from 'react';
import { shallow } from 'enzyme';
import List from './index';

function defaultProps() {
  return {
    funds: [],
    filterBy: () => null,
  };
}

it('should always render as a div', () => {
  const props = defaultProps();
  const wrapper = shallow(<List {...props} />);
  expect(wrapper.type()).toEqual('div');
});

it('should render at least one header row', () => {
  const props = defaultProps();
  const wrapper = shallow(<List {...props} />);
  expect(wrapper.find('Row').length).toEqual(1);
});

it('should render one row per funds', () => {
  const props = defaultProps();
  props.funds = [
    { id: '1000', rating: 1, '1m': 0.0, '3m': 0.0, '6m': 0.0, '1y': 0.0, v3y: 0.0, score: 0.0 },
    { id: '2000', rating: 1, '1m': 0.0, '3m': 0.0, '6m': 0.0, '1y': 0.0, v3y: 0.0, score: 0.0 },
    { id: '8000', rating: 1, '1m': 0.0, '3m': 0.0, '6m': 0.0, '1y': 0.0, v3y: 0.0, score: 0.0 },
  ];
  const wrapper = shallow(<List {...props} />);

  expect(wrapper.find('Row').length).toEqual(4);
});
