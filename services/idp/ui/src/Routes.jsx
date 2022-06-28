import React, { lazy } from 'react';
import PropTypes from 'prop-types';

import { Route, Switch } from 'react-router-dom';

import PrivateRoute from './components/PrivateRoute';

const AsyncLogin = lazy(() =>
  import(/* webpackChunkName: "containers-login" */ './containers/Login'));
const AsyncWelcome = lazy(() =>
  import(/* webpackChunkName: "containers-welcome" */ './containers/Welcome'));
const AsyncGoodbye = lazy(() =>
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
