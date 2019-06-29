import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import Button from 'components/Button';
import { COLUMNS } from 'components/Funds/Constants';
import style from './index.module.css';

export default function HeaderIcon({ filter, onClick, icon, displayed }) {
  const list = Object.entries(COLUMNS)
    .filter(([, column]) => column[filter])
    .map(([key, column]) => (
      <li key={key}>
        <Button type="none" className={style.button} onClick={() => onClick(key)}>
          {column.label}
        </Button>
      </li>
    ));

  const iconClasses = classnames({
    [style.active]: displayed,
  })

  const listClasses = classnames({
    [style.displayed]: displayed,
    [style.hidden]: !displayed,
  })

  return (
    <span className={style.icon}>
      <span className={iconClasses}>{icon}</span>
      <ol className={listClasses}>{list}</ol>
    </span>
  );
}

HeaderIcon.displayName = 'HeaderIcon';

HeaderIcon.propTypes = {
  filter: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  icon: PropTypes.node.isRequired,
  displayed: PropTypes.bool.isRequired,
};
