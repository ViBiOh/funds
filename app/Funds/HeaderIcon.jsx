import React from 'react';
import PropTypes from 'prop-types';
import style from './HeaderIcon.css';

const HeaderIcon = ({ columns, filter, onClick, icon, displayed }) => {
  const list = Object.keys(columns).filter(e => columns[e][filter]).map(key => (
    <li key={key}>
      <button onClick={() => onClick(key)}>{columns[key].label}</button>
    </li>
  ));

  return (
    <span className={style.icon}>
      <span className={displayed ? style.active : ''}>
        {icon}
      </span>
      <ol className={displayed ? style.displayed : style.hidden}>
        {list}
      </ol>
    </span>
  );
};

HeaderIcon.displayName = 'HeaderIcon';

HeaderIcon.propTypes = {
  columns: PropTypes.shape({}).isRequired,
  filter: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  icon: PropTypes.node.isRequired,
  displayed: PropTypes.bool.isRequired,
};

export default HeaderIcon;
