import React from 'react';
import { makeRouteConfig, Route } from 'found';
import { graphql } from 'react-relay';

import App from './components/App';
import EventsList from './components/EventsList';
import ChecksList from './components/CheckList';

const AppQuery = graphql`
  query routes_App_Query {
    viewer {
      ...App_viewer
    }
  }
`;

const ListQuery = graphql`
  query routes_EventsList_Query {
    viewer {
      ...EventsList_viewer
    }
  }
`;

const CheckRouteQuery = graphql`
  query routes_CheckList_Query {
    viewer {
      ...CheckList_viewer
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
      path="events"
      Component={EventsList}
      query={ListQuery}
    />
    <Route
      path="checks"
      Component={ChecksList}
      query={CheckRouteQuery}
    />
  </Route>,
);
