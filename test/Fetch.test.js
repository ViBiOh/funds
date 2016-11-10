/* eslint-disable import/no-extraneous-dependencies */
/* eslint-env mocha */
import { expect } from 'chai';
import Fetch from '../app/Fetch';

describe('Fetch', () => {
  it('should return a promise', () => {
    global.fetch = () => Promise.resolve({
      status: 200,
      headers: {
        get: () => '',
      },
      text: () => Promise.resolve(''),
    });

    const result = Fetch.get('/');

    expect(result).to.be.instanceof(Promise);
  });

  it('should return text when asked', () => {
    global.fetch = () => Promise.resolve({
      status: 200,
      headers: {
        get: () => 'text/plain',
      },
      text: () => Promise.resolve('Test Mocha'),
    });

    return Fetch.get('/').then(data => expect(data).to.eql('Test Mocha'));
  });

  it('should return json when asked', () => {
    global.fetch = () => Promise.resolve({
      status: 200,
      headers: {
        get: () => 'application/json',
      },
      json: () => Promise.resolve({
        result: 'Test Mocha',
      }),
    });

    return Fetch.get('/').then(data => expect(data).to.eql({
      result: 'Test Mocha',
    }));
  });

  it('should return error when 400 or more', () => {
    global.fetch = () => Promise.resolve({
      status: 400,
      headers: {
        get: () => 'text/plain',
      },
      text: () => Promise.resolve('Test Mocha Error'),
    });

    return Fetch.get('/').catch(data => expect(data.content).to.eql('Test Mocha Error'));
  });

  it('should return jsonError when 400 or more', () => {
    global.fetch = () => Promise.resolve({
      status: 500,
      headers: {
        get: () => 'application/json',
      },
      json: () => Promise.resolve({
        error: 'Test Mocha Error',
      }),
    });

    return Fetch.get('/').catch(data => expect(data.content).to.eql({
      error: 'Test Mocha Error',
    }));
  });

  it('should return error when json fail', () => {
    global.fetch = () => Promise.resolve({
      status: 200,
      headers: {
        get: () => 'application/json',
      },
      json: () => Promise.reject(new Error('Mocha JSON Error')),
    });

    return Fetch.get('/').catch(data => expect(String(data)).to.eql('Error: Mocha JSON Error'));
  });
});
