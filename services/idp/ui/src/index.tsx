import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

import './i18n';

import App from './App';
import store from './store';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store as any}>
      <App/>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
