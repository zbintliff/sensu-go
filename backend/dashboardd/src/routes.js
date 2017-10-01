import React from 'react';
import { makeRouteConfig, Route } from 'found';
import { graphql } from 'react-relay';

import App from './components/App';
import EventsList from './components/eventsList';

const AppQuery = graphql`
  query routes_App_Query {
    viewer {
      ...App_viewer
    }
  }
`;

export default makeRouteConfig(
  <Route
    path="/"
    Component={App}
    query={AppQuery}
  >,
    <Route
      Component={EventsList}
    />
  </Route>,
);
