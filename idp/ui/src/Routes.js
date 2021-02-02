import React from 'react';
import PropTypes from 'prop-types';

import { Route, Switch } from 'react-router-dom';
import AsyncComponent from 'kpop/es/AsyncComponent';

import PrivateRoute from './components/PrivateRoute';

const AsyncLogin = AsyncComponent(() =>
  import(/* webpackChunkName: "containers-login" */ './containers/Login'));
const AsyncWelcome = AsyncComponent(() =>
  import(/* webpackChunkName: "containers-welcome" */ './containers/Welcome'));
const AsyncGoodbye = AsyncComponent(() =>
  import(/* webpackChunkName: "containers-goodbye" */ './containers/Goodbye'));

const Routes = ({ hello }) => (
  <Switch>
    <PrivateRoute
      path="/welcome"
      exact
      component={AsyncWelcome}
      hello={hello}
    />
    <Route
      path="/goodbye"
      exact
      component={AsyncGoodbye}
    />
    <Route
      path="/"
      component={AsyncLogin}
    />
  </Switch>
);

Routes.propTypes = {
  hello: PropTypes.object
};

export default Routes;
