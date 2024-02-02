import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import './i18n';

import App from './App';
import store from './store';

import './app.css';

const root = document.getElementById('root')

// if a custom background image has been configured, make use of it
const bgImg = root.getAttribute('data-bg-img')

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store as any}>
      <App bgImg={bgImg}/>
    </Provider>
  </React.StrictMode>,
  root
);
