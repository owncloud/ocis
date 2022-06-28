import React, { Suspense, lazy } from 'react';

import { MuiThemeProvider } from '@material-ui/core/styles';
import {
  CssBaseline,
 } from '@material-ui/core';

import './App.css';
import './fancy-background.css';
import Spinner from './components/Spinner';
import * as version from './version';
import theme from './theme';

const LazyMain = lazy(() => import(/* webpackChunkName: "identifier-main" */ './Main'));

console.info(`Kopano Identifier build version: ${version.build}`); // eslint-disable-line no-console

const App = () => {
  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline/>
      <Suspense fallback={<Spinner/>}>
        <LazyMain />
      </Suspense>
    </MuiThemeProvider>
  );
}

export default App;
