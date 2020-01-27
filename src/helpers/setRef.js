/**
 * Define ref for component.
 * @param {ReactComponent} that      React component to put ref in
 * @param {String}         name      Ref name
 * @param {ReactComponent} component Component to assign to the ref
 */
export default function setRef(that, name, component) {
  that[name] = component; // eslint-disable-line no-param-reassign
}
