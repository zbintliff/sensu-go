import React from 'react';
import PropTypes from 'prop-types';
import { TableRow, TableCell } from 'material-ui/Table';
import Checkbox from 'material-ui/Checkbox';
import { createFragmentContainer, graphql } from 'react-relay';

class EventRow extends React.Component {
  render() {
    const { event: { entity, config, timestamp }, ...other } = this.props;
    return (
      <TableRow {...other}>
        <TableCell checkbox><Checkbox /></TableCell>
        <TableCell>{entity.entityID}</TableCell>
        <TableCell>{config.name}</TableCell>
        <TableCell>{config.command}</TableCell>
        <TableCell>{timestamp}</TableCell>
      </TableRow>
    );
  }
}

EventRow.propTypes = {
  event: PropTypes.shape({
    entity: PropTypes.shape({ id: '' }).isRequired,
    config: PropTypes.shape({ name: '', command: '' }).isRequired,
    timestamp: PropTypes.string.isRequired,
  }).isRequired,
};

export default createFragmentContainer(
  EventRow,
  graphql`
    fragment EventRow_event on Event {
      ... on CheckEvent {
        timestamp
        config {
          name
          command
        }
        entity {
          entityID
        }
      }
    }
  `,
);
