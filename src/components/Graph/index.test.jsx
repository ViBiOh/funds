import React from 'react';
import { shallow } from 'enzyme';
import Graph from './index';

jest.mock('chart.js', () => jest.fn().mockImplementation(() => ({
  data: {
    datasets: [],
    labels: [],
  },
  update: jest.fn(),
  destroy: jest.fn(),
})));

function defaultProps() {
  return {
    data: {
      datasets: [],
      labels: [],
    },
  };
}

it('should always render as a canvas', () => {
  const props = defaultProps();
  const wrapper = shallow(<Graph {...props} />);
  expect(wrapper.type()).toEqual('canvas');
});

it('should do nothing if chart not here on unmount', () => {
  const props = defaultProps();
  const wrapper = shallow(<Graph {...props} />);
  wrapper.unmount();
});
