/**
 * Insert Fathom script into DOM.
 */
function insertScript(url) {
  return new Promise((resolve) => {
    const script = document.createElement('script');
    script.id = 'fathom-script';
    script.type = 'text/javascript';
    script.src = `${url}/tracker.js`;
    script.async = 'true';
    script.onload = resolve;

    document.querySelector('head').appendChild(script);
  });
}

/**
 * Initialize Fathom analytics tracking.
 * @param  {Object} config Configuration ojbect
 */
function init(config) {
  insertScript(config.FATHOM_URL).then(() => {
    global.fathom('trackPageview');
  });
}

export default {
  init,
  insertScript,
};
