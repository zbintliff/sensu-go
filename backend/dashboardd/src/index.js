import React from 'react';
import ReactDOM from 'react-dom';

import Routes from './components/Routes';
import registerServiceWorker from './registerServiceWorker';

import './index.css';

// Renderer
ReactDOM.render(
  <Routes />,
  document.getElementById('root'),
);

// Register service workers
registerServiceWorker();
