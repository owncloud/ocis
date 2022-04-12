import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import DialogActions from '@material-ui/core/DialogActions';

import ResponsiveScreen from '../../components/ResponsiveScreen';
import { executeLogoff } from '../../actions/common';

const styles = theme => ({
  button: {
    margin: theme.spacing(1),
    minWidth: 100
  },
  subHeader: {
    marginBottom: theme.spacing(5)
  }
});

class Welcomescreen extends React.PureComponent {
  render() {
    const { classes, hello } = this.props;

    const loading = hello === null;
    return (
      <ResponsiveScreen loading={loading}>
        <Typography variant="h5" component="h3" className="oc-light">
          <FormattedMessage
            id="konnect.welcome.headline"
            defaultMessage="Welcome {displayName}"
            values={{displayName: hello.displayName}}>
          </FormattedMessage>
        </Typography>
        <Typography variant="subtitle1" className={classes.subHeader + " oc-light"}>
          {hello.username}
        </Typography>

        <Typography gutterBottom className="oc-light">
          <FormattedMessage id="konnect.welcome.message"
            defaultMessage="You are signed in - awesome!"></FormattedMessage>
        </Typography>

        <DialogActions>
          <Button
            color="secondary"
            className={classes.button}
            variant="contained"
            onClick={(event) => this.logoff(event)}
          >
            <FormattedMessage id="konnect.welcome.signoutButton.label" defaultMessage="Sign out"></FormattedMessage>
          </Button>
        </DialogActions>
      </ResponsiveScreen>
    );
  }

  logoff(event) {
    event.preventDefault();

    this.props.dispatch(executeLogoff()).then((response) => {
      const { history } = this.props;

      if (response.success) {
        history.push('/identifier');
      }
    });
  }
}

Welcomescreen.propTypes = {
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

export default connect(mapStateToProps)(withStyles(styles)(Welcomescreen));
