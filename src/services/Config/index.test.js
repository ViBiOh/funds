import funtch from 'funtch';
import FundsService from './index';

jest.mock('funtch');

it('should fetch config from env', () => {
  return FundsService.getFunds().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('/env');
  });
});
