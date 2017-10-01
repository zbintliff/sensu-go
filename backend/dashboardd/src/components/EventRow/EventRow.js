import React from 'react';
import PropTypes from 'prop-types';
import { TableRow, TableRowColumn } from 'material-ui/Table';
import { createFragmentContainer, graphql } from 'react-relay';

class EventRow extends React.Component {
  render() {
    const { event: { entity, config, timestamp }, ...other } = this.props;
    return (
      <TableRow {...other}>
        {other.children[0] /* checkbox passed down from TableBody */}
        <TableRowColumn>{entity.id}</TableRowColumn>
        <TableRowColumn>{config.name}</TableRowColumn>
        <TableRowColumn>{config.command}</TableRowColumn>
        <TableRowColumn>{timestamp}</TableRowColumn>
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
          id
        }
      }
    }
  `,
);
