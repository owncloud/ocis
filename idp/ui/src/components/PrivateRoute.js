import React from 'react';
import PropTypes from 'prop-types';
import { Route } from 'react-router-dom';

import RedirectWithQuery from './RedirectWithQuery';

const PrivateRoute = ({ component: Target, hello, ...rest }) => (
  <Route {...rest} render={props => (
    hello ? (
      <Target {...props}/>
    ) : (
      <RedirectWithQuery target='/identifier' />
    )
  )}/>
);

PrivateRoute.propTypes = {
  component: PropTypes.func.isRequired,
  hello: PropTypes.object
};

export default PrivateRoute;
