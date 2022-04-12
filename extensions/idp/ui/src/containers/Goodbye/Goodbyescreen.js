import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import renderIf from 'render-if';
import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import DialogActions from '@material-ui/core/DialogActions';

import ResponsiveScreen from '../../components/ResponsiveScreen';
import { executeHello, executeLogoff } from '../../actions/common';

const styles = theme => ({
  subHeader: {
    marginBottom: theme.spacing(5)
  },
  wrapper: {
    marginTop: theme.spacing(5),
    position: 'relative',
    display: 'inline-block'
  }
});

class Goodbyescreen extends React.PureComponent {
  componentDidMount() {
    this.props.dispatch(executeHello());
  }

  render() {
    const { classes, hello } = this.props;

    const loading = hello === null;
    return (
      <ResponsiveScreen loading={loading}>
        {renderIf(hello !== null && !hello.state)(() => (
          <div>
            <Typography variant="h5" component="h3">
              <FormattedMessage id="konnect.goodbye.headline" defaultMessage="Goodbye"></FormattedMessage>
            </Typography>
            <Typography variant="subtitle1" className={classes.subHeader}>
              <FormattedMessage id="konnect.goodbye.subHeader"
                defaultMessage="you have been signed out from your Kopano account">
              </FormattedMessage>
            </Typography>
            <Typography gutterBottom>
              <FormattedMessage id="konnect.goodbye.message.close"
                defaultMessage="You can close this window now.">
              </FormattedMessage>
            </Typography>
          </div>
        ))}
        {renderIf(hello !== null && hello.state === true)(() => (
          <div>
            <Typography variant="h5" component="h3">
              <FormattedMessage
                id="konnect.goodbye.confirm.headline"
                defaultMessage="Hello {displayName}"
                values={{displayName: hello.displayName}}>
              </FormattedMessage>
            </Typography>
            <Typography variant="subtitle1" className={classes.subHeader}>
              <FormattedMessage id="konnect.goodbye.confirm.subHeader"
                defaultMessage="please confirm sign out">
              </FormattedMessage>
            </Typography>

            <Typography gutterBottom>
              <FormattedMessage id="konnect.goodbye.message.confirm"
                defaultMessage="Press the button below, to sign out from your Kopano account now.">
              </FormattedMessage>
            </Typography>

            <DialogActions>
              <div className={classes.wrapper}>
                <Button
                  color="secondary"
                  className={classes.button}
                  onClick={(event) => this.logoff(event)}
                >
                  <FormattedMessage id="konnect.goodbye.signoutButton.label"
                    defaultMessage="Sign out"></FormattedMessage>
                </Button>
              </div>
            </DialogActions>
          </div>
        ))}
      </ResponsiveScreen>
    );
  }

  logoff(event) {
    event.preventDefault();

    this.props.dispatch(executeLogoff()).then((response) => {
      const { history } = this.props;

      if (response.success) {
        this.props.dispatch(executeHello());
        history.push('/goodbye');
      }
    });
  }
}

Goodbyescreen.propTypes = {
  classes: PropTypes.object.isRequired,

  hello: PropTypes.object,

  dispatch: PropTypes.func.isRequired,
  history: PropTypes.object.isRequired
};

const mapStateToProps = (state) => {
  const { hello } = state.common;

  return {
    hello
  };
};

export default connect(mapStateToProps)(withStyles(styles)(Goodbyescreen));
