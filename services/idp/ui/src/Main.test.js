import React from 'react';
import ReactDOM from 'react-dom';

import { Provider } from 'react-redux';

import store from './store';
import App from './Main';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Provider store={store}><App/></Provider>, div);
});
