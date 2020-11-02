import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import green from '@material-ui/core/colors/green';
import Typography from '@material-ui/core/Typography';

import TextInput from '../../components/TextInput'

import { updateInput, executeLogonIfFormValid, advanceLogonFlow } from '../../actions/login';
import { ErrorMessage } from '../../errors';

const styles = theme => ({
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12
  },
  subHeader: {
    marginBottom: theme.spacing(3)
  },
  wrapper: {
    position: 'relative',
    width: '100%',
    textAlign: 'center'
  },
  message: {
    marginTop: 5,
    marginBottom: 5
  }
});

class Login extends React.PureComponent {
  state = {};

  componentDidMount() {
    const { hello, query, dispatch, history } = this.props;
    if (hello && hello.state && history.action !== 'PUSH') {
      if (!query.prompt || query.prompt.indexOf('select_account') == -1) {
        dispatch(advanceLogonFlow(true, history));
        return;
      }

      history.replace(`/chooseaccount${history.location.search}${history.location.hash}`);
      return;
    }
  }

  render() {
    const { loading, errors, classes, username } = this.props;
    const hasError = errors.http || errors.username || errors.password;
    const errorMessage = errors.http
      ? <ErrorMessage error={errors.http}></ErrorMessage>
      : (errors.username
        ? <ErrorMessage error={errors.username}></ErrorMessage>
        : <ErrorMessage error={errors.password}></ErrorMessage>);

    return (
      <form action="" onSubmit={(event) => this.logon(event)}>
        <TextInput
          autoFocus
          autoCapitalize="off"
          spellCheck="false"
          value={username}
          onChange={this.handleChange('username')}
          autoComplete="kopano-account username"
          placeholder={({ id: "konnect.login.usernameField.label", defaultMessage: "Username" })}
        />
        <TextInput
          type="password"
          margin="normal"
          onChange={this.handleChange('password')}
          autoComplete="kopano-account current-password"
          placeholder={({ id: "konnect.login.usernameField.label", defaultMessage: "Password" })}
        />
        {hasError && <Typography variant="subtitle2" color="error" className={classes.message}>{errorMessage}</Typography>}
        <div className={classes.wrapper}>
          <Button
            type="submit"
            color="primary"
            variant="contained"
            className="oc-button-primary oc-mt-l"
            disabled={!!loading}
            onClick={(event) => this.logon(event)}
          >
            <FormattedMessage id="konnect.login.nextButton.label" defaultMessage="Log in"></FormattedMessage>
          </Button>
          {loading && <CircularProgress size={24} className={classes.buttonProgress} />}
        </div>
      </form>
    );
  }

  handleChange(name) {
    return event => {
      this.props.dispatch(updateInput(name, event.target.value));
    };
  }

  logon(event) {
    event.preventDefault();

    const { username, password, dispatch, history } = this.props;
    dispatch(executeLogonIfFormValid(username, password, false)).then((response) => {
      if (response.success) {
        dispatch(advanceLogonFlow(response.success, history));
      }
    });
  }
}

Login.propTypes = {
  classes: PropTypes.object.isRequired,

  loading: PropTypes.string.isRequired,
  username: PropTypes.string.isRequired,
  password: PropTypes.string.isRequired,
  errors: PropTypes.object.isRequired,
  hello: PropTypes.object,
  query: PropTypes.object.isRequired,

  dispatch: PropTypes.func.isRequired,
  history: PropTypes.object.isRequired
};

const mapStateToProps = (state) => {
  const { loading, username, password, errors} = state.login;
  const { hello, query } = state.common;

  return {
    loading,
    username,
    password,
    errors,
    hello,
    query
  };
};

export default connect(mapStateToProps)(withStyles(styles)(Login));
