import funtch from 'funtch';
import ConfigService from './index';

jest.mock('funtch');

it('should fetch config from env', () => {
  return ConfigService.getConfig().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('/env');
  });
});
