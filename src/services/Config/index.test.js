import funtch from 'funtch';
import store from 'AppStore';
import ConfigService from './index';

jest.mock('funtch');
jest.mock('AppStore');

it('should fetch config from env', () => {
  return ConfigService.getConfig().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('/env');
  });
});

it('should handle null config from store', () => {
  store.getState.mockReturnValue({});

  expect(ConfigService.getApiUrl()).toEqual('');
});

it('should concat given path to API url', () => {
  store.getState.mockReturnValue({
    config: {
      API_URL: 'https://api.vibioh.fr',
    },
  });

  expect(ConfigService.getApiUrl('list')).toEqual('https://api.vibioh.fr/list');
});
