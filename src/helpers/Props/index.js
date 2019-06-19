import PropTypes from 'prop-types';

/**
 * String or number for a prop type.
 * @type {Object}
 */
export const STRING_OR_NUMBER = PropTypes.oneOfType([PropTypes.string, PropTypes.number]);
