import React from 'react';
import PropTypes from 'prop-types';
import map from 'lodash/map';

import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
} from 'material-ui/Table';
import EventRow from '../eventRow';

const styles = require('./eventsList.css');

function EventsList({ events }) {
  console.info('dfasdfasdf', events);
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
        {map(events, (event, i) => (
          <EventRow key={i} {...event} />
        ))}
      </TableBody>
    </Table>
  );
}

EventsList.propTypes = {
  events: PropTypes.arrayOf(PropTypes.object).isRequired,
};

export default EventsList;
