import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import { createLogger } from 'redux-logger';

import rootReducer from './reducers';

const middlewares = [
  thunkMiddleware
];

if (process.env.NODE_ENV === 'development') { // eslint-disable-line no-undef
  middlewares.push(createLogger()); // must be last middleware in the chain.
}

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
  rootReducer,
  composeEnhancers(applyMiddleware(
    ...middlewares,
  ))
);

export default store;
