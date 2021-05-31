import React, { Component } from 'react';
import PropTypes from 'prop-types';
import {
  Chart,
  BarElement,
  BarController,
  LinearScale,
  CategoryScale,
  Legend,
  Tooltip,
} from 'chart.js';
import setRef from 'helpers/setRef';

export default class Graph extends Component {
  /**
   * React lifecycle.
   */
  componentDidMount() {
    this.updateChart(this.props);
  }

  /**
   * React lifecycle.
   */
  componentDidUpdate() {
    this.updateChart(this.props);
  }

  /**
   * React lifecycle.
   */
  componentWillUnmount() {
    this.clearChart();
  }

  /**
   * Update chart value.
   * @param  {Object} config Chart configuration
   * @return {Object}        Chart value
   */
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
          x: {
            display: false,
          },
          y: {
            display: false,
            beginAtZero: true,
          },
        },
      };

      Chart.register(
        BarElement,
        BarController,
        LinearScale,
        CategoryScale,
        Legend,
        Tooltip,
      );
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
    return (
      <canvas
        ref={(e) => setRef(this, 'graph', e)}
        className={className}
        height={200}
      />
    );
  }
}

Graph.propTypes = {
  className: PropTypes.string,
  data: PropTypes.shape({
    datasets: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
    labels: PropTypes.arrayOf(PropTypes.string).isRequired,
  }).isRequired,
};

Graph.defaultProps = {
  className: '',
};
