import { createStore, applyMiddleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import appReducers from 'reducers';
import appSagas from 'sagas';

/**
 * ReduxSaga configuration.
 * @type {Object}
 */
const sagaMiddleware = createSagaMiddleware();

/**
 * Redux store.
 * @type {Object}
 */
const appStore = createStore(appReducers, applyMiddleware(sagaMiddleware));

sagaMiddleware.run(appSagas);

/**
 * Redux's store of application.
 */
export default appStore;
