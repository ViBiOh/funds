import funtch from 'funtch';
import { getApiUrl } from 'Constants';
import FundsService from './index';

jest.mock('funtch');
jest.mock('Constants');

it('should fetch data for funds', () => {
  getApiUrl.mockReturnValue('localhost');
  funtch.get.mockResolvedValue({});

  return FundsService.getFunds().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('localhost/list');
  });
});
