import React from 'react';
import PropTypes from 'prop-types';
import { CHART_COLORS } from 'components/Funds/Constants';
import ChartGraph from 'components/Graph';
import style from './index.module.css';

export default function Graph({ aggregat, aggregated }) {
  if (!aggregat.key) {
    return null;
  }

  const data = {
    labels: [],
    datasets: [
      {
        label: 'Aggregat count ',
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

  return <ChartGraph type="bar" data={data} className={style.container} />;
}

Graph.displayName = 'Graph';

Graph.propTypes = {
  aggregat: PropTypes.shape({
    key: PropTypes.string,
  }).isRequired,
  aggregated: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string.isRequired,
      count: PropTypes.number.isRequired,
    }),
  ).isRequired,
};
