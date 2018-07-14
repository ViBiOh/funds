import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Chart from 'chart.js';
import setRef from '../Tools/ref';

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
    const { type, data } = config;

    if (this.chart) {
      this.chart.data.datasets = data.datasets;
      this.chart.data.labels = data.labels;
      this.chart.update();
    } else if (this.graph) {
      const options = {
        legend: {
          display: false,
        },
        scales: {
          xAxes: [
            {
              display: false,
            },
          ],
          yAxes: [
            {
              display: false,
              ticks: {
                beginAtZero: true,
              },
            },
          ],
        },
      };

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
    const { className } = this.props;
    return <canvas ref={e => setRef(this, 'graph', e)} className={className} height={200} />;
  }
}

Graph.propTypes = {
  className: PropTypes.string,
};

Graph.defaultProps = {
  className: '',
};
