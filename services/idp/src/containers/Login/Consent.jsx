import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { withTranslation, Trans } from 'react-i18next';

import renderIf from 'render-if';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import BaseTooltip from '@material-ui/core/Tooltip';
import CircularProgress from '@material-ui/core/CircularProgress';
import green from '@material-ui/core/colors/green';
import Typography from '@material-ui/core/Typography';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';

import { executeConsent, advanceLogonFlow, receiveValidateLogon } from '../../actions/login';
import { ErrorMessage } from '../../errors';
import { REQUEST_CONSENT_ALLOW } from '../../actions/types';
import ClientDisplayName from '../../components/ClientDisplayName';
import ScopesList from '../../components/ScopesList';

const styles = theme => ({
  button: {
    margin: theme.spacing(1),
    minWidth: 100
  },
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12
  },
  subHeader: {
    marginBottom: theme.spacing(2)
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

const Tooltip = ({children, ...other } = {}) => {
  // Ensures that there is only a single child for the tooltip element to
  // make it compatible with the Trans component.
  return <BaseTooltip {...other}><span>{children}</span></BaseTooltip>;
}

Tooltip.propTypes = {
    children: PropTypes.node,
};

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
    const { classes, loading, hello, errors, client, t } = this.props;

    const scopes = hello.details.scopes || {};
    const meta = hello.details.meta || {};

    return (
      <DialogContent>
        <Typography variant="h5" component="h3" className="oc-light">
          {t("konnect.consent.headline", "Hi {{displayName}}", { displayName: hello.displayName })}
        </Typography>
        <Typography variant="subtitle1" className={classes.subHeader + " oc-light oc-mb-m"}>
          {hello.username}
        </Typography>

        <Typography variant="subtitle1" gutterBottom className="oc-light">
          <Trans t={t} i18nKey="konnect.consent.message">
            <Tooltip
              placement="bottom"
              title={t("konnect.consent.tooltip.client", 'Clicking "Allow" will redirect you to: {{redirectURI}}', { redirectURI: client.redirect_uri })}
            >
              <em><ClientDisplayName client={client}/></em>
            </Tooltip> wants to
          </Trans>
        </Typography>
        <ScopesList dense disablePadding className={classes.scopesList} scopes={scopes} meta={meta.scopes}></ScopesList>

        <Typography variant="subtitle1" gutterBottom className="oc-light">
          <Trans t={t} i18nKey="konnect.consent.question">
            Allow <em><ClientDisplayName client={client}/></em> to do this?
          </Trans>
        </Typography>
        <Typography className="oc-light">
          {t("konnect.consent.consequence", "By clicking Allow, you allow this app to use your information.")}
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
                {t("konnect.consent.cancelButton.label", "Cancel")}
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
                {t("konnect.consent.allowButton.label", "Allow")}
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
      </DialogContent>
    );
  }
}

Consent.propTypes = {
  classes: PropTypes.object.isRequired,
  t: PropTypes.func.isRequired,

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

export default connect(mapStateToProps)(withStyles(styles)(withTranslation()(Consent)));
