import 'kpop/static/css/base.css';
import 'kpop/static/css/scrollbar.css';
import 'typeface-roboto';
import './app.css';
import './fancy-background.css';

import * as kpop from 'kpop/es/version';

import * as version from './version';

console.info(`Kopano Identifier build version: ${version.build}`); // eslint-disable-line no-console
console.info(`Kopano Kpop build version: ${kpop.build}`); // eslint-disable-line no-console

// NOTE(longsleep): Load async, this enables code splitting via Webpack.
import(/* webpackChunkName: "identifier-app" */ './app');
