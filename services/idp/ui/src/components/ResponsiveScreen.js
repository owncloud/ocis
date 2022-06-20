import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { FormattedMessage } from 'react-intl';

import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import DialogContent from '@material-ui/core/DialogContent';

import Loading from './Loading';

const styles = theme => ({
  root: {
    display: 'flex',
    flex: 1
  },
  content: {
    position: 'relative',
    width: '100%'
  },
  actions: {
    marginTop: -40,
    justifyContent: 'flex-start',
    paddingLeft: theme.spacing(3),
    paddingRight: theme.spacing(3)
  },
  wrapper: {
    width: '100%',
    maxWidth: 300,
    display: 'flex',
    flex: 1,
    alignItems: 'center'
  }
});

const footerProductName = name => <strong>{name}</strong>;

const ResponsiveScreen = (props) => {
  const {
    classes,
    withoutLogo,
    withoutPadding,
    loading,
    children,
    className,
    DialogProps,
    PaperProps,
    ...other
  } = props;

  const logo = withoutLogo ? null :
    <img src={process.env.PUBLIC_URL + '/static/logo.svg'} className="oc-logo" alt="ownCloud Logo"/>;

  const content = loading ? <Loading/> : (withoutPadding ? children : <DialogContent>{children}</DialogContent>);

  return (
    <Grid container justify="center" alignItems="center" direction="column" spacing={0}
      className={classNames(classes.root, className)} {...other}>
      <div className={classes.wrapper}>
        <div className={classes.content}>
          {logo}
          {content}
        </div>
      </div>
      <footer className="oc-footer-message">
        <FormattedMessage
          id="konnect.footer.slogan"
          defaultMessage="<name>ownCloud</name> - a safe home for all your data"
          values={{
            name: chunks => footerProductName(chunks)
          }}
        />
      </footer>
    </Grid>
  );
};

ResponsiveScreen.defaultProps = {
  withoutLogo: false,
  withoutPadding: false,
  loading: false
};

ResponsiveScreen.propTypes = {
  classes: PropTypes.object.isRequired,
  withoutLogo: PropTypes.bool,
  withoutPadding: PropTypes.bool,
  loading: PropTypes.bool,
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  PaperProps: PropTypes.object,
  DialogProps: PropTypes.object
};

export default withStyles(styles)(ResponsiveScreen);
