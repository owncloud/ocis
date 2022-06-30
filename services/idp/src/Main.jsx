import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { BrowserRouter } from 'react-router-dom';

import Routes from './Routes';

class Main extends PureComponent {
  render() {
    const { hello, pathPrefix } = this.props;

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

Main.propTypes = {
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

export default connect(mapStateToProps)(Main);
