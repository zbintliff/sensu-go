import React from 'react';
import {
  createBrowserRouter,
  // HttpError,
  makeRouteConfig,
  // Redirect,
  Route,
} from 'found';

import AppContainer from './App';

const Router = createBrowserRouter({
  routeConfig: makeRouteConfig(
    <Route
      path="/"
      Component={AppContainer}
    />,
  ),

  renderError: ({ error }) => ( // eslint-disable-line
    <div>
      {error.status === 404 ? 'Not found' : 'Error'}
    </div>
  ),
});

export default Router;
