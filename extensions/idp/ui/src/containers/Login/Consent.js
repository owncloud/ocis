import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import renderIf from 'render-if';
import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Tooltip from '@material-ui/core/Tooltip';
import CircularProgress from '@material-ui/core/CircularProgress';
import green from '@material-ui/core/colors/green';
import Typography from '@material-ui/core/Typography';
import DialogActions from '@material-ui/core/DialogActions';

import { executeConsent, advanceLogonFlow, receiveValidateLogon } from '../../actions/login';
import { ErrorMessage } from '../../errors';
import { REQUEST_CONSENT_ALLOW } from '../../actions/types';
import ClientDisplayName from '../../components/ClientDisplayName';
import ScopesList from '../../components/ScopesList';

const styles = theme => ({
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12
  },
  scopesList: {
    marginBottom: theme.spacing(2)
  },
  wrapper: {
    marginTop: theme.spacing(2),
    position: 'relative',
    display: 'inline-block'
  },
  message: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2)
  }
});

class Consent extends React.PureComponent {
  componentDidMount() {
    const { dispatch, hello, history, client } = this.props;
    if ((!hello || !hello.state || !client) && history.action !== 'PUSH') {
      history.replace(`/identifier${history.location.search}${history.location.hash}`);
    }

    dispatch(receiveValidateLogon({})); // XXX(longsleep): hack to reset loading and errors.
  }

  action = (allow=false, scopes={}) => (event) => {
    event.preventDefault();

    if (allow === undefined) {
      return;
    }

    // Convert all scopes which are true to a scope value.
    const scope = Object.keys(scopes).filter(scope => {
      return !!scopes[scope];
    }).join(' ');

    const { dispatch, history } = this.props;
    dispatch(executeConsent(allow, scope)).then((response) => {
      if (response.success) {
        dispatch(advanceLogonFlow(response.success, history, true, {konnect: response.state}));
      }
    });
  }

  render() {
    const { classes, loading, hello, errors, client } = this.props;

    const scopes = hello.details.scopes || {};
    const meta = hello.details.meta || {};

    return (
      <div>
        <Typography variant="h5" component="h3" className="oc-light">
          <FormattedMessage
            id="konnect.consent.headline"
            defaultMessage="Hi {displayName}"
            values={{displayName: hello.displayName}}
          />
        </Typography>
        <Typography variant="subtitle1" className="oc-light oc-mb-m">
          {hello.username}
        </Typography>

        <Typography variant="subtitle1" gutterBottom className="oc-light">
          <FormattedMessage
            id="konnect.consent.message"
            defaultMessage="{clientDisplayName} wants to"
            values={{clientDisplayName:
              <Tooltip
                placement="bottom"
                title={<FormattedMessage
                  id="konnect.consent.tooltip.client"
                  defaultMessage='Clicking "Allow" will redirect you to: {redirectURI}'
                  values={{
                    redirectURI: client.redirect_uri
                  }}
                ></FormattedMessage>}
              >
                <em><ClientDisplayName client={client}/></em>
              </Tooltip>
            }}
          ></FormattedMessage>
        </Typography>
        <ScopesList dense disablePadding className={classes.scopesList} scopes={scopes} meta={meta.scopes}></ScopesList>

        <Typography className="oc-light">
          <FormattedMessage
            id="konnect.consent.consequence"
            defaultMessage="By clicking Allow, you allow this app to use your information.">
          </FormattedMessage>
        </Typography>

        <form action="" onSubmit={this.action(undefined, scopes)}>
          <DialogActions>
            <div className={classes.wrapper}>
              <Button
                color="secondary"
                className={classes.button}
                disabled={!!loading}
                onClick={this.action(false, scopes)}
              >
                <FormattedMessage id="konnect.consent.cancelButton.label" defaultMessage="Cancel"></FormattedMessage>
              </Button>
              {(loading && loading !== REQUEST_CONSENT_ALLOW) &&
                <CircularProgress size={24} className={classes.buttonProgress} />}
            </div>
            <div className={classes.wrapper}>
              <Button
                type="submit"
                color="primary"
                variant="contained"
                className="oc-button-primary"
                disabled={!!loading}
                onClick={this.action(true, scopes)}
              >
                <FormattedMessage id="konnect.consent.allowButton.label" defaultMessage="Allow"></FormattedMessage>
              </Button>
              {loading === REQUEST_CONSENT_ALLOW && <CircularProgress size={24} className={classes.buttonProgress} />}
            </div>
          </DialogActions>

          {renderIf(errors.http)(() => (
            <Typography variant="subtitle2" color="error" className={classes.message}>
              <ErrorMessage error={errors.http}></ErrorMessage>
            </Typography>
          ))}
        </form>
      </div>
    );
  }
}

Consent.propTypes = {
  classes: PropTypes.object.isRequired,

  loading: PropTypes.string.isRequired,
  errors: PropTypes.object.isRequired,
  hello: PropTypes.object,
  client: PropTypes.object.isRequired,

  dispatch: PropTypes.func.isRequired,
  history: PropTypes.object.isRequired
};

const mapStateToProps = (state) => {
  const { hello } = state.common;
  const { loading, errors } = state.login;

  return {
    loading: loading,
    errors,
    hello,
    client: hello.details.client || {}
  };
};

export default connect(mapStateToProps)(withStyles(styles)(Consent));
