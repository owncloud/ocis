import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { BrowserRouter } from 'react-router-dom';

import Routes from './Routes';

class App extends PureComponent {
  render() {
    const { classes, hello, pathPrefix } = this.props;

    return (
      <BrowserRouter basename={pathPrefix}>
        <Routes hello={hello}/>
      </BrowserRouter>
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

export default connect(mapStateToProps)(App);
