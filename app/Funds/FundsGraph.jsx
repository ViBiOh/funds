import React from 'react';
import PropTypes from 'prop-types';
import { CHART_COLORS } from './FundsConstantes';
import Graph from './Graph';
import style from './FundsGraph.less';

const FundsGraph = ({ aggregat, aggregated }) => {
  if (!aggregat.key) {
    return null;
  }

  const options = {
    legend: false,
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

  const data = {
    labels: [],
    datasets: [
      {
        label: 'Count',
        data: [],
        backgroundColor: [],
      },
    ],
  };

  let i = 0;
  aggregated.forEach((entry) => {
    data.labels.push(entry.label);
    data.datasets[0].data.push(entry.count);
    data.datasets[0].backgroundColor.push(CHART_COLORS[i]);

    i = (i + 1) % CHART_COLORS.length;
  });

  return <Graph type="bar" data={data} options={options} className={style.container} />;
};

FundsGraph.displayName = 'FundsGraph';

FundsGraph.propTypes = {
  aggregat: PropTypes.shape({}).isRequired,
  aggregated: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
};

export default FundsGraph;
