import React, { Component } from 'react';
import Chartist from 'chartist';

export default class Graph extends Component {
  componentDidMount() {
    this.updateChart(this.props);
  }

  componentWillReceiveProps(newProps) {
    this.updateChart(newProps);
  }

  componentWillUnmount() {
    if (this.chartist) {
      try {
        this.chartist.detach();
      } catch (err) {
        throw new Error('Internal chartist error', err);
      }
    }
  }

  updateChart(config) {
    const { type, data, options, responsive } = config;

    if (this.chartist) {
      this.chartist.update(data, options, responsive);
    } else {
      this.chartist = new Chartist[type](this.chart, data, options, responsive);
    }

    return this.chartist;
  }

  render() {
    const className = `ct-chart ${this.props.className}`;

    return (
      <div className={className} key="graph" ref={e => (this.chart = e)} />
    );
  }
}

Graph.propTypes = {
  className: React.PropTypes.string,
};
