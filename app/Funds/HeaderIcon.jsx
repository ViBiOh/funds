import React from 'react';
import PropTypes from 'prop-types';
import Button from '../Button/Button';
import { COLUMNS } from './FundsConstantes';
import style from './HeaderIcon.less';

const HeaderIcon = ({ filter, onClick, icon, displayed }) => {
  const list = Object.keys(COLUMNS)
    .filter(e => COLUMNS[e][filter])
    .map(key => (
      <li key={key}>
        <Button type="none" className={style.button} onClick={() => onClick(key)}>
          {COLUMNS[key].label}
        </Button>
      </li>
    ));

  return (
    <span className={style.icon}>
      <span className={displayed ? style.active : ''}>{icon}</span>
      <ol className={displayed ? style.displayed : style.hidden}>{list}</ol>
    </span>
  );
};

HeaderIcon.displayName = 'HeaderIcon';

HeaderIcon.propTypes = {
  filter: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  icon: PropTypes.node.isRequired,
  displayed: PropTypes.bool.isRequired,
};

export default HeaderIcon;
