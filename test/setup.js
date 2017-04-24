/* eslint-disable import/no-extraneous-dependencies */
import { jsdom } from 'jsdom';

global.document = new JSDOM('<html><body></body></html>').window.document;
global.window = document.defaultView;
global.navigator = window.navigator;
