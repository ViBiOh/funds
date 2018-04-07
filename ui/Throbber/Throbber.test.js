import test from 'ava';
import React from 'react';
import { shallow } from 'enzyme';
import Throbber from './Throbber';

test('should render into a div', (t) => {
  t.is(shallow(<Throbber />).type(), 'div');
});

test('should have no label by default', (t) => {
  t.is(shallow(<Throbber />).find('span').length, 0);
});

test('should have label when given', (t) => {
  t.is(
    shallow(<Throbber label="test" />)
      .find('span')
      .text(),
    'test',
  );
});