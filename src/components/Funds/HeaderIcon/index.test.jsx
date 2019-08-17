import React from 'react';
import { shallow } from 'enzyme';
import { COLUMNS } from 'components/Funds/Constants';
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

it('should call given callback on item click', () => {
  const props = defaultProps();
  props.onClick = jest.fn();

  const wrapper = shallow(<HeaderIcon {...props} />);
  wrapper
    .find('Button')
    .at(0)
    .simulate('click');

  expect(props.onClick).toHaveBeenCalledWith(Object.keys(COLUMNS)[0]);
});
