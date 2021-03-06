/**
 * Action creator : return a function for given action
 * @param  {string}     type     Action type
 * @param  {...objects} argNames Properties' names of action
 * @return {func}                Function that generate action with type and properties given the
 *                               params
 */
export const makeActionCreator = (type, ...argNames) => (...args) => {
  const action = { type };
  argNames.forEach((arg, index) => {
    action[argNames[index]] = args[index];
  });

  return action;
};

/**
 * Transform a name into a type name : SNAKE_UPPER_CASE
 * @param  {string} name A camel case action name
 * @return {string}      Snake upper case type name
 */
export const toTypeName = (name) => String(name)
  .replace(/([A-Z])/g, '_$1')
  .toUpperCase();

/**
 * Action creator : return the function and the constant for the given action
 * @param  {string} type   Action type
 * @param  {string} action Action function name
 * @param  {Array}  inputs Properties' names of action
 * @return {object}        An object containing both function and constant
 */
export const makeActionAndTypeCreator = (type, action, inputs = []) => ({
  [type]: type,
  [action]: makeActionCreator(type, ...inputs),
});

/**
 * Action creator for an API call (request, success, fail)
 * @param  {string} camelCaseName CamelCase name of action : the action function name
 * @param  {Array}  inputs        Properties' names of request action
 * @param  {Array}  outputs       Properties' names of response action
 * @return {object}               An object container constants and functions for requesting API
 */
export const makeApiActionCreator = (
  camelCaseName,
  inputs = [],
  outputs = [],
) => {
  const typeName = toTypeName(camelCaseName);

  return {
    [typeName]: typeName,
    ...makeActionAndTypeCreator(`${typeName}_REQUEST`, camelCaseName, inputs),
    ...makeActionAndTypeCreator(
      `${typeName}_SUCCEEDED`,
      `${camelCaseName}Succeeded`,
      outputs,
    ),
    ...makeActionAndTypeCreator(
      `${typeName}_FAILED`,
      `${camelCaseName}Failed`,
      ['error'],
    ),
  };
};
