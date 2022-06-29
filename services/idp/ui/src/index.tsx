import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import './i18n';

import App from './App';
import store from './store';

import './app.css';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store as any}>
      <App/>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
