import PropTypes from 'prop-types';

/**
 * String or number for a prop type.
 * @type {Object}
 */
// eslint-disable-next-line import/prefer-default-export
export const STRING_OR_NUMBER = PropTypes.oneOfType([
  PropTypes.string,
  PropTypes.number,
]);
