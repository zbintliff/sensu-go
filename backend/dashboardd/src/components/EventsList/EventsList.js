import React from 'react';
import PropTypes from 'prop-types';
import map from 'lodash/map';
import { createFragmentContainer, graphql } from 'react-relay';

import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
} from 'material-ui/Table';
import EventRow from '../EventRow';

const styles = require('./eventsList.css');

class EventsList extends React.Component {
  static propTypes = {
    viewer: PropTypes.shape({
      events: PropTypes.array,
    }).isRequired,
  }

  render() {
    const { viewer } = this.props;

    return (
      <Table className={styles.table}>
        <TableHeader>
          <TableRow>
            <TableHeaderColumn>Entity</TableHeaderColumn>
            <TableHeaderColumn>Check</TableHeaderColumn>
            <TableHeaderColumn>Command</TableHeaderColumn>
            <TableHeaderColumn>Timestamp</TableHeaderColumn>
          </TableRow>
        </TableHeader>
        <TableBody>
          {map(viewer.events, (event, i) => (
            <EventRow key={i} event={event} />
          ))}
        </TableBody>
      </Table>
    );
  }
}


export default createFragmentContainer(
  EventsList,
  graphql`
    fragment EventsList_viewer on Viewer {
      events {
        ...EventRow_event
      }
    }
  `,
);
