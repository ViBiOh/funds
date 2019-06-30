import React from 'react';
import { shallow } from 'enzyme';
import Header from './index';

function defaultProps() {
  return {
    aggregateBy: () => null,
    filterBy: () => null,
    orderBy: () => null,
  };
}

it('should always render as a header', () => {
  const props = defaultProps();
  const wrapper = shallow(<Header {...props} />);
  expect(wrapper.type()).toEqual('header');
});
