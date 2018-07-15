/**
 * Insert Fathom script into DOM.
 */
function insertScript(url) {
  const trackerScript = `${url}/tracker.js`;

  return new Promise((resolve) => {
    global.fathom = global.fathom
      || function initQ() {
        (global.fathom.q = global.fathom.q || []).push(document, global, trackerScript, 'fathom');
      };

    const script = document.createElement('script');
    script.id = 'fathom-script';
    script.type = 'text/javascript';
    script.src = trackerScript;
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

/**
 * Track page view.
 */
function track() {
  if (global.fathom) {
    global.fathom('trackPageview');
  }
}

export default {
  init,
  track,
};
