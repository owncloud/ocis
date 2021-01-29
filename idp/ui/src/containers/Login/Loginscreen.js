import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { Route, Switch } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';

import ResponsiveScreen from '../../components/ResponsiveScreen';
import RedirectWithQuery from '../../components/RedirectWithQuery';
import { executeHello } from '../../actions/common';

import Login from './Login';
import Chooseaccount from './Chooseaccount';
import Consent from './Consent';

const styles = () => ({
});

class Loginscreen extends React.PureComponent {
  componentDidMount() {
    this.props.dispatch(executeHello());
  }

  render() {
    const { hello } = this.props;

    const loading = hello === null;
    return (
      <ResponsiveScreen loading={loading} withoutPadding={true} >
        <Switch>
          <Route path="/identifier" exact component={Login}></Route>
          <Route path="/chooseaccount" exact component={Chooseaccount}></Route>
          <Route path="/consent" exact component={Consent}></Route>
          <RedirectWithQuery target="/identifier"/>
        </Switch>
      </ResponsiveScreen>
    );
  }
}

Loginscreen.propTypes = {
  classes: PropTypes.object.isRequired,

  hello: PropTypes.object,

  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = (state) => {
  const { hello } = state.common;

  return {
    hello
  };
};

export default connect(mapStateToProps)(withStyles(styles)(Loginscreen));
