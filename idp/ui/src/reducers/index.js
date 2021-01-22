import { combineReducers } from 'redux';

import commonReducer from './common';
import loginReducer from './login';

const rootReducer = combineReducers({
  common: commonReducer,
  login: loginReducer
});

export default rootReducer;
