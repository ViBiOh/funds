import test from 'ava';
import sinon from 'sinon';
import funtch from 'funtch';
import FundsService from './FundsService';

test.beforeEach(() => {
  sinon.stub(funtch, 'get').callsFake(url => Promise.resolve({ url }));
});

test.afterEach(() => {
  funtch.get.restore();
});

test('should fetch data for funds', t => FundsService.getFunds().then(({ url }) => {
  t.truthy(/list/i.test(url));
}));
