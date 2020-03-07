import React from 'react';
import sinon from 'sinon';
import { shallow } from 'enzyme';
import ConnectedFunds from 'containers/Funds';
import { AppComponent } from './index';

function defaultProps() {
  return {
    ready: false,
    init: () => null,
  };
}

it('should render as a div if not ready', () => {
  const props = defaultProps();
  const wrapper = shallow(<AppComponent {...props} />);
  expect(wrapper.type()).toEqual('div');
});

it('should show a Throbber while not ready', () => {
  const props = defaultProps();

  const wrapper = shallow(<AppComponent {...props} />);
  expect(wrapper.find('Throbber').length).toEqual(1);
});

it('should trigger init on mount', () => {
  const props = defaultProps();
  props.init = sinon.spy();

  const wrapper = shallow(<AppComponent {...props} />);
  wrapper.instance().componentDidMount();

  expect(props.init.called).toEqual(true);
});

it('should render Funds when ready', () => {
  const props = defaultProps();
  props.ready = true;

  const wrapper = shallow(<AppComponent {...props} />);
  expect(wrapper.find(ConnectedFunds).length).toEqual(1);
});
