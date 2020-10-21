import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router';
import { Redirect } from 'react-router-dom';

const RedirectWithQuery = ({target, location, ...rest}) => {
  const to = {
    pathname: target,
    search: location.search,
    hash: location.hash
  };

  return (
    <Redirect to={to} {...rest}></Redirect>
  );
};

RedirectWithQuery.propTypes = {
  target: PropTypes.string.isRequired,
  location: PropTypes.object.isRequired
};

export default withRouter(RedirectWithQuery);
