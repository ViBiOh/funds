import React from 'react';
import { shallow } from 'enzyme';
import Graph from './index';

function defaultProps() {
  return {};
}

it('should always render as a canvas', () => {
  const props = defaultProps();
  const wrapper = shallow(<Graph {...props} />);
  expect(wrapper.type()).toEqual('canvas');
});
