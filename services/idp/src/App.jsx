import React, { Suspense, lazy } from 'react';

import { MuiThemeProvider } from '@material-ui/core/styles';
import { defaultTheme as theme } from 'kpop/es/theme';

import 'kpop/static/css/base.css';
import 'kpop/static/css/scrollbar.css';
import './App.css';

import Spinner from './components/Spinner';
import * as version from './version';

const LazyMain = lazy(() => import(/* webpackChunkName: "identifier-main" */ './Main'));

console.info(`Kopano Identifier build version: ${version.build}`); // eslint-disable-line no-console

const App = () => {
  return (
    <MuiThemeProvider theme={theme}>
      <Suspense fallback={<Spinner/>}>
        <LazyMain />
      </Suspense>
    </MuiThemeProvider>
  );
}

export default App;
