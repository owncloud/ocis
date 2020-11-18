import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { FormattedMessage } from 'react-intl';

import LinearProgress from '@material-ui/core/LinearProgress';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import renderIf from 'render-if';

import { retryHello } from '../actions/common';
import { ErrorMessage } from '../errors';

function Loading({ error, dispatch }) {
  const retry = (event) => {
    event.preventDefault();
    dispatch(retryHello());
  }

  return (
    <Grid item align="center">
      {renderIf(error === null)(() => (
        <LinearProgress className="oc-progress" />
      ))}
      {renderIf(error !== null)(() => (
        <div>
          <Typography className="oc-light" variant="h5" gutterBottom align="center">
            <FormattedMessage id="konnect.loading.error.headline" defaultMessage="Failed to connect to server" />
          </Typography>
          <Typography align="center" color="error">
            <ErrorMessage error={error} />
          </Typography>
          <Button
            autoFocus
            color="primary"
            variant="contained"
            className="oc-button-primary oc-mt-l"
            onClick={(event) => retry(event)}
          >
            <FormattedMessage id="konnect.login.retryButton.label" defaultMessage="Retry" />
          </Button>
        </div>
      ))}
    </Grid>
  );
}

Loading.propTypes = {
  error: PropTypes.object,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = (state) => {
  const { error } = state.common;

  return {
    error
  };
};

export default connect(mapStateToProps)(Loading);
