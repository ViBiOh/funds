import test from 'ava';
import sinon from 'sinon';
import funtch from 'funtch';
import FundsService from './FundsService';
import { MS_URL } from '../Constants';

test.beforeEach(() => {
  sinon.stub(funtch, 'get').callsFake(url => Promise.resolve({ url }));
});

test.afterEach(() => {
  funtch.get.restore();
});

test('should get data source url for given id', (t) => {
  t.is(FundsService.getDataUrl('test'), MS_URL('test'));
});

test('should fetch data for funds', t =>
  FundsService.getFunds().then(({ url }) => {
    t.truthy(/list/i.test(url));
  }));
