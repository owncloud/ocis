import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import Avatar from '@material-ui/core/Avatar';
import Typography from '@material-ui/core/Typography';

import { executeLogonIfFormValid, advanceLogonFlow } from '../../actions/login';
import { ErrorMessage } from '../../errors';

const styles = theme => ({
  subHeader: {
    marginBottom: theme.spacing(2)
  },
  message: {
    marginTop: theme.spacing(2)
  },
  accountList: {
    marginLeft: theme.spacing(-3),
    marginRight: theme.spacing(-3)
  },
  accountListItem: {
    paddingLeft: theme.spacing(3),
    paddingRight: theme.spacing(3)
  }
});

function Chooseaccount({ loading, errors, classes, hello, history, dispatch }) {
  useEffect(() => {
    if ((!hello || !hello.state) && history.action !== 'PUSH') {
      history.replace(`/identifier${history.location.search}${history.location.hash}`);
    }
  });

  let errorMessage = null;

  if (errors.http) {
    errorMessage = <Typography color="error" className={classes.message}>
      <ErrorMessage error={errors.http}></ErrorMessage>
    </Typography>;
  }

  let username = '';

  if (hello && hello.state) {
    username = hello.username;
  }

  const logon = (event) => {
    event.preventDefault();
    dispatch(executeLogonIfFormValid(hello.username, '', true)).then((response) => {
      if (response.success) {
        dispatch(advanceLogonFlow(response.success, history));
      }
    });
  }

  const logoff = (event) => {
    event.preventDefault();
    history.push(`/identifier${history.location.search}${history.location.hash}`);
  }

  return (
    <div>
      <Typography variant="h5" component="h3" className="oc-light">
        <FormattedMessage id="konnect.chooseaccount.headline" defaultMessage="Choose an account"></FormattedMessage>
      </Typography>
      <Typography variant="subtitle1" className={classes.subHeader + " oc-light"}>
        <FormattedMessage id="konnect.chooseaccount.subHeader" defaultMessage="to sign in to Kopano">
        </FormattedMessage>
      </Typography>

      <form action="" onSubmit={(event) => logon(event)}>
        <List disablePadding className={classes.accountList}>
          <ListItem
            button
            disableGutters
            className={classes.accountListItem}
            disabled={!!loading}
            onClick={(event) => logon(event)}
          ><ListItemAvatar><Avatar>{username.substr(0, 1)}</Avatar></ListItemAvatar>
            <ListItemText className="oc-light" primary={username} />
          </ListItem>
          <ListItem
            button
            disableGutters
            className={classes.accountListItem}
            disabled={!!loading}
            onClick={(event) => logoff(event)}
          >
            <ListItemAvatar>
              <Avatar>
                <FormattedMessage id="konnect.chooseaccount.useOther.persona.label" defaultMessage="?">
                </FormattedMessage>
              </Avatar>
            </ListItemAvatar>
            <ListItemText
              className="oc-light"
              primary={
                <FormattedMessage
                  id="konnect.chooseaccount.useOther.label"
                  defaultMessage="Use another account">
                </FormattedMessage>
              }
            />
          </ListItem>
        </List>

        {errorMessage}
      </form>
    </div>
  );
}

Chooseaccount.propTypes = {
  classes: PropTypes.object.isRequired,

  loading: PropTypes.string.isRequired,
  errors: PropTypes.object.isRequired,
  hello: PropTypes.object,

  dispatch: PropTypes.func.isRequired,
  history: PropTypes.object.isRequired
};

const mapStateToProps = (state) => {
  const { loading, errors } = state.login;
  const { hello } = state.common;

  return {
    loading,
    errors,
    hello
  };
};

export default connect(mapStateToProps)(withStyles(styles)(Chooseaccount));
