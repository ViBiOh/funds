import React, { Component } from 'react';
import Chart from 'chart.js';

export default class Graph extends Component {
  componentDidMount() {
    this.updateChart(this.props);
  }

  componentWillReceiveProps(newProps) {
    this.updateChart(newProps);
  }

  componentWillUnmount() {
    this.clearChart();
  }

  updateChart(config) {
    const { type, data, options } = config;

    if (this.chart) {
      this.chart.data.datasets = data.datasets;
      this.chart.data.labels = data.labels;
      this.chart.update();
    } else {
      this.chart = new Chart(this.graph, {
        type,
        data,
        options,
      });
    }

    return this.chart;
  }

  clearChart() {
    if (this.chart) {
      this.chart.destroy();
    }
  }

  render() {
    return (
      <canvas
        ref={e => (this.graph = e)}
        className={this.props.className}
        height={200}
      />
    );
  }
}

Graph.propTypes = {
  className: React.PropTypes.string,
};

Graph.defaultProps = {
  className: '',
};

