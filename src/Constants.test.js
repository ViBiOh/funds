import funtch from 'funtch';
import Constants from './Constants';

jest.mock('funtch');

test('should fetch data from /env', () => {
  funtch.get.mockResolvedValue({ API_URL: 'localhost' });

  return Constants.init().then(() => {
    expect(funtch.get).toHaveBeenCalledWith('/env');
  })
});

test('should return API_URL from context', () => {
  funtch.get.mockResolvedValue({ API_URL: 'localhost' });

  return Constants.init().then(() => {
    expect(Constants.getApiUrl()).toEqual('localhost');
  });
});
