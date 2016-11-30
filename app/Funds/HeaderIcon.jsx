import React from 'react';
import style from './Funds.css';

const HeaderIcon = ({ columns, filter, onClick, icon, displayed }) => {
  const list = Object.keys(columns)
    .filter(e => columns[e][filter])
    .map(key => (
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
  columns: React.PropTypes.shape({}).isRequired,
  filter: React.PropTypes.string.isRequired,
  onClick: React.PropTypes.func.isRequired,
  icon: React.PropTypes.node.isRequired,
  displayed: React.PropTypes.bool.isRequired,
};

export default HeaderIcon;
