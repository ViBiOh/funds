import funtch from 'funtch';
import ConfigService from 'services/Config';
import FundsService from './index';

jest.mock('funtch');
jest.mock('services/Config');

it('should fetch data for funds', () => {
  ConfigService.getApiUrl.mockReturnValue('localhost');
  funtch.get.mockResolvedValue({});

  return FundsService.getFunds().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('localhost');
  });
});
