import React from 'react';
import { shallow } from 'enzyme';
import Throbber from './index';

test('should render into a div', () => {
  expect(shallow(<Throbber />).type()).toEqual('div');
});

test('should have no label by default', () => {
  expect(shallow(<Throbber />).find('span').length).toEqual(0);
});

test('should have label when given', () => {
  expect(
    shallow(<Throbber label="test" />)
      .find('span')
      .text(),
  ).toEqual('test');
});
