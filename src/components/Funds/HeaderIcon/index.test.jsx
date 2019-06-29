import React from 'react';
import { shallow } from 'enzyme';
import HeaderIcon from './index';

function defaultProps() {
  return {
    filter: 'sortable',
    onClick: () => null,
    icon: <span />,
    displayed: false,
  };
}

it('should always render as a span', () => {
  const props = defaultProps();
  const wrapper = shallow(<HeaderIcon {...props} />);
  expect(wrapper.type()).toEqual('span');
});

