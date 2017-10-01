import {
  Environment,
  Network,
  RecordSource,
  Store,
} from 'relay-runtime';

function fetchQuery(
  operation,
  variables,
  // cacheConfig,
  // uploadables,
) {
  // TODO: Make URL configurable
  return fetch('//localhost:8080/graphql', {
    method: 'POST',
    headers: {
      // Add authentication and other headers here
      'content-type': 'application/json',
    },
    body: JSON.stringify({
      query: operation.text, // GraphQL text from input
      variables,
    }),
  }).then(response => response.json());
}

// Create a record source & instantiate store
const source = new RecordSource();
const store = new Store(source);

// Create a network layer from the fetch function
const network = Network.create(fetchQuery);

// Create an environment using this network:
const environment = new Environment({
  network,
  source,
  store,
});

export default environment;
