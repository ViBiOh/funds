import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import actions from 'actions';
import ConnectedFunds from 'containers/Funds';
import Throbber from 'components/Throbber';
import style from './index.module.css';

/**
 * App Component.
 */
export class AppComponent extends Component {
  /**
   * React lifecycle.
   */
  componentDidMount() {
    const { init } = this.props;
    init();
  }

  /**
   * React lifecycle.
   */
  render() {
    const { ready } = this.props;

    if (!ready) {
      return (
        <div className={style.loader}>
          <Throbber label="Loading environment" />
        </div>
      );
    }

    return <ConnectedFunds />;
  }
}

AppComponent.propTypes = {
  ready: PropTypes.bool.isRequired,
  init: PropTypes.func.isRequired,
};

/**
 * Select props from Redux state.
 * @param {Object} state Current state
 */
function mapStateToProps(state) {
  return {
    ready: state.config.ready,
  };
}

/**
 * Provide actions to dispatch.
 * @type {Object}
 */
const mapDispatchToProps = {
  init: actions.init,
};

/**
 * AppComponent connected.
 */
export default connect(mapStateToProps, mapDispatchToProps)(AppComponent);
