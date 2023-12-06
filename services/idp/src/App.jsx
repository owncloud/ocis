import React, { ReactElement, Suspense, lazy } from 'react';
import PropTypes from 'prop-types';

import { MuiThemeProvider } from '@material-ui/core/styles';
import { defaultTheme as theme } from 'kpop/es/theme';

import 'kpop/static/css/base.css';
import 'kpop/static/css/scrollbar.css';

import Spinner from './components/Spinner';
import * as version from './version';

const LazyMain = lazy(() => import(/* webpackChunkName: "identifier-main" */ './Main'));

console.info(`Kopano Identifier build version: ${version.build}`); // eslint-disable-line no-console

const App = ({ bgImg }): ReactElement => {
  return (
    <div
      className='oc-login-bg'
      style={{ backgroundImage: bgImg ? `url(${bgImg})` : undefined }}
    >
      <MuiThemeProvider theme={theme}>
        <Suspense fallback={<Spinner/>}>
          <LazyMain/>
        </Suspense>
      </MuiThemeProvider>
    </div>
  );
}

App.propTypes = {
  bgImg: PropTypes.string
};

export default App;
