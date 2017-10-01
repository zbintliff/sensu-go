import React from 'react';
import ReactDOM from 'react-dom';
import { Resolver } from 'found-relay';
import BrowserProtocol from 'farce/lib/BrowserProtocol';
import queryMiddleware from 'farce/lib/queryMiddleware';
import createFarceRouter from 'found/lib/createFarceRouter';
import createRender from 'found/lib/createRender';

import routes from './routes';
import registerServiceWorker from './registerServiceWorker';
import environment from './environment';

import './index.css';

const Router = createFarceRouter({
  historyProtocol: new BrowserProtocol(),
  historyMiddlewares: [queryMiddleware],
  routeConfig: routes,
  render: createRender({}),
});

// Renderer
ReactDOM.render(
  <Router resolver={new Resolver(environment)} />,
  document.getElementById('root'),
);

// Register service workers
registerServiceWorker();
