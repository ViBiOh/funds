#!/usr/bin/env node

/* eslint-disable no-console */

/**
 * We want to convert files with hashes (in path and in content) to files without hashes
 * in order to avoid breaking changes.
 *
 * During deployment
 *  * client can receive `index.html` from oldversion with ref to `index.abcde1234.js`,
 *  * client ask for `index.abcde1234.js`, but because of load-balancing, fall to newVersion
 *    who don't know it.
 *
 * To avoid this, we remove every hash in file
 * (keeping files with hashes in case if we miss something):
 *   * in filename
 *   * in content of files
 */
const fs = require('fs').promises;

const OUTPUT_DIR = 'build';
const HASH_FINDER = /((static\/[^")\\]*?)\.([a-z0-9]{8})(\.[^")\\]+))/gim;

const processedFiles = [];
const excludePattern = RegExp('/media/', 'gim');
const version = process.argv[2];

if (!version) {
  console.error('You must provide version hash for cache-busting');
  process.exit(1);
}

console.info(`Using '${version}' as cache-buster version`);

/**
 * Replace and remove hashes URI in file
 * @param {string} file         Filename to clean
 * @param {string} hashlessFile Output filename
 */
async function replaceHash(file, hashlessFile) {
  processedFiles.push(file);

  excludePattern.lastIndex = 0;
  if (excludePattern.test(file)) {
    console.info(`Copying file from ${file} to ${hashlessFile || file}`);
    await fs.copyFile(file, hashlessFile);
    return;
  }

  console.info(`Converting hashes in ${file} into ${hashlessFile || file}`);

  let fileContent = await fs.readFile(file, 'utf-8');
  fileContent = fileContent.replace(
    HASH_FINDER,
    (all, filename, prefix, hash, suffix) => {
      const buildFilename = `${OUTPUT_DIR}/${filename}`;
      const outputFilename = `${prefix}${suffix}`;

      if (processedFiles.indexOf(buildFilename) === -1) {
        replaceHash(buildFilename, `${OUTPUT_DIR}/${outputFilename}`);
      }

      return `${outputFilename}?v=${version}`;
    },
  );

  if (hashlessFile) {
    await fs.writeFile(hashlessFile, fileContent);
  } else {
    await fs.writeFile(file, fileContent);
  }
}

replaceHash(`${OUTPUT_DIR}/index.html`);
replaceHash(`${OUTPUT_DIR}/asset-manifest.json`);
