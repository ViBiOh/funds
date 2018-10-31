import React from 'react';
import { shallow } from 'enzyme';
import Button from '.';

test('should always render as a button', () => {
  expect(shallow(<Button />).type()).toEqual('button');
});

test('should not wrap child', () => {
  const wrapper = shallow(
    <Button>
      <span>
First
      </span>
    </Button>,
  );

  expect(wrapper.find('span').length).toEqual(1);
});

test('should wrap children in div', () => {
  const wrapper = shallow(
    <Button>
      <span>
First
      </span>
      <span>
Second
      </span>
    </Button>,
  );

  expect(wrapper.find('div').length).toEqual(1);
});
