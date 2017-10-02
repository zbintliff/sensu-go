import React from 'react';
import PropTypes from 'prop-types';
import map from 'lodash/map';
import { createFragmentContainer, graphql } from 'react-relay';

import Table, {
  TableBody,
  TableHead,
  TableCell,
  TableRow,
} from 'material-ui/Table';
import Checkbox from 'material-ui/Checkbox';
import Row from './CheckRow';

const styles = require('./List.css');

class CheckList extends React.Component {
  static propTypes = {
    viewer: PropTypes.shape({
      checks: PropTypes.shape({
        edges: PropTypes.array.isRequired,
      }),
    }).isRequired,
  }

  render() {
    const { viewer } = this.props;

    return (
      <Table className={styles.table}>
        <TableHead>
          <TableRow>
            <TableCell checkbox><Checkbox /></TableCell>
            <TableCell>Check</TableCell>
            <TableCell>Command</TableCell>
            <TableCell>Subscribers</TableCell>
            <TableCell>Interval</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {map(viewer.checks.edges, edge =>
            <Row key={edge.cursor} check={edge.node} />,
          )}
        </TableBody>
      </Table>
    );
  }
}


export default createFragmentContainer(
  CheckList,
  graphql`
    fragment CheckList_viewer on Viewer {
      checks(first: 200) {
        edges {
          node {
            ...CheckRow_check
          }
          cursor
        }
        pageInfo {
          hasNextPage
        }
      }
    }
  `,
);
