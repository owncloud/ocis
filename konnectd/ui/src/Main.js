import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { BrowserRouter } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';

import { enhanceBodyBackground } from './utils';
import Routes from './Routes';

// Trigger loading of background image.
enhanceBodyBackground();

const styles = () => ({
  root: {
    position: 'relative',
    display: 'flex',
    flex: 1
  }
});

class App extends PureComponent {
  render() {
    const { classes, hello, pathPrefix } = this.props;

    return (
      <div className={classes.root}>
        <BrowserRouter basename={pathPrefix}>
          <Routes hello={hello}/>
        </BrowserRouter>
      </div>
    );
  }

  reload(event) {
    event.preventDefault();

    window.location.reload();
  }
}

App.propTypes = {
  classes: PropTypes.object.isRequired,

  hello: PropTypes.object,
  updateAvailable: PropTypes.bool.isRequired,
  pathPrefix: PropTypes.string.isRequired
};

const mapStateToProps = (state) => {
  const { hello, updateAvailable, pathPrefix } = state.common;

  return {
    hello,
    updateAvailable,
    pathPrefix
  };
};

export default connect(mapStateToProps)(withStyles(styles)(App));
