import React from 'react';
import { shallow } from 'enzyme';
import Graph from './index';

function defaultProps() {
  return {
    aggregat: {},
    aggregated: [],
  };
}

it('should not render if no aggregat key', () => {
  const props = defaultProps();
  const wrapper = shallow(<Graph {...props} />);
  expect(wrapper.type()).toEqual(null);
});

it('should render a Graph if key and aggregated values', () => {
  const props = defaultProps();
  props.aggregat.key = 'test';
  props.aggregated = [
    {
      label: 'first',
      count: 8,
    },
    {
      label: 'second',
      count: 12,
    },
  ];

  const wrapper = shallow(<Graph {...props} />);
  expect(wrapper.type().name).toEqual('Graph');
});
