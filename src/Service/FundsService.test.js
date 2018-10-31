import funtch from 'funtch';
import Constants from '../Constants';
import FundsService from './FundsService';

jest.mock('funtch');
jest.mock('../Constants');

test('should fetch data for funds', () => {
  Constants.getApiUrl.mockReturnValue('localhost');
  funtch.get.mockResolvedValue({});

  return FundsService.getFunds().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('localhost/list');
  });
});
