import test from 'ava';
import React from 'react';
import { shallow } from 'enzyme';
import PerformanceCell from './PerformanceCell';

const defaultProps = {
  type: '1y',
};

test('should always render as a span', (t) => {
  t.is(shallow(<PerformanceCell {...defaultProps} />).type(), 'span');
});

test('should have style for given value', (t) => {
  const wrapper = shallow(<PerformanceCell {...defaultProps} value={10} />);

  t.is(wrapper.props().className.split(' ').length, 3);
});

test('should have style for negative value', (t) => {
  const wrapper = shallow(<PerformanceCell {...defaultProps} value={-10} />);

  t.is(wrapper.props().className.split(' ').length, 3);
});
