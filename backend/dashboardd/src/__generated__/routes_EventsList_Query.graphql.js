/**
 * @flow
 * @relayHash f516c1b3f9310694c7b69eda3c029827
 */

/* eslint-disable */

'use strict';

/*::
import type {ConcreteBatch} from 'relay-runtime';
export type routes_EventsList_QueryResponse = {|
  +viewer: ?{| |};
|};
*/


/*
query routes_EventsList_Query {
  viewer {
    ...EventsList_viewer
  }
}

fragment EventsList_viewer on Viewer {
  events {
    __typename
    ...EventRow_event
    ... on Node {
      id
    }
    ... on MetricEvent {
      id
    }
  }
}

fragment EventRow_event on Event {
  ... on CheckEvent {
    timestamp
    config {
      name
      command
    }
    entity {
      id
    }
  }
}
*/

const batch /*: ConcreteBatch*/ = {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "routes_EventsList_Query",
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "viewer",
        "args": null,
        "concreteType": "Viewer",
        "name": "__viewer_viewer",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "EventsList_viewer",
            "args": null
          }
        ],
        "storageKey": null
      }
    ],
    "type": "Query"
  },
  "id": null,
  "kind": "Batch",
  "metadata": {},
  "name": "routes_EventsList_Query",
  "query": {
    "argumentDefinitions": [],
    "kind": "Root",
    "name": "routes_EventsList_Query",
    "operation": "query",
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "args": null,
        "concreteType": "Viewer",
        "name": "viewer",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "args": null,
            "concreteType": null,
            "name": "events",
            "plural": true,
            "selections": [
              {
                "kind": "ScalarField",
                "alias": null,
                "args": null,
                "name": "__typename",
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "args": null,
                "name": "id",
                "storageKey": null
              },
              {
                "kind": "InlineFragment",
                "type": "CheckEvent",
                "selections": [
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "args": null,
                    "name": "timestamp",
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "args": null,
                    "concreteType": "CheckConfig",
                    "name": "config",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "args": null,
                        "name": "name",
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "args": null,
                        "name": "command",
                        "storageKey": null
                      }
                    ],
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "args": null,
                    "concreteType": "Entity",
                    "name": "entity",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "args": null,
                        "name": "id",
                        "storageKey": null
                      }
                    ],
                    "storageKey": null
                  }
                ]
              }
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "args": null,
        "handle": "viewer",
        "name": "viewer",
        "key": "",
        "filters": null
      }
    ]
  },
  "text": "query routes_EventsList_Query {\n  viewer {\n    ...EventsList_viewer\n  }\n}\n\nfragment EventsList_viewer on Viewer {\n  events {\n    __typename\n    ...EventRow_event\n    ... on Node {\n      id\n    }\n    ... on MetricEvent {\n      id\n    }\n  }\n}\n\nfragment EventRow_event on Event {\n  ... on CheckEvent {\n    timestamp\n    config {\n      name\n      command\n    }\n    entity {\n      id\n    }\n  }\n}\n"
};

module.exports = batch;
