import React from 'react';
import { shallow, mount } from 'enzyme';
import Chart from 'chart.js';
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

it('should update chart on mount', () => {
  const props = defaultProps();
  mount(<Graph {...props} />);
  expect(Chart).toHaveBeenCalled();
});

it('should update chart on update', () => {
  const props = defaultProps();
  const wrapper = mount(<Graph {...props} />);
  wrapper.setProps({ type: 'bar' });

  expect(wrapper.instance().chart.update).toHaveBeenCalled();
});

it('should destroy chart on unmount', () => {
  const props = defaultProps();
  const wrapper = mount(<Graph {...props} />);
  const destroyFn = wrapper.instance().chart.destroy;
  wrapper.unmount();

  expect(destroyFn).toHaveBeenCalled();
});
