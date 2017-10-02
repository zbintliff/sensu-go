import React from 'react';
import PropTypes from 'prop-types';
import { MuiThemeProvider, createMuiTheme } from 'material-ui/styles';
import { createFragmentContainer, graphql } from 'react-relay';

import Sidebar from '../sidebar';
import Toolbar from '../toolbar';

const muiTheme = createMuiTheme({
  palette: {
    primary1Color: '#92C72E',
  },
});

const styles = require('./app.css');

class App extends React.Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
  };

  state = {
    toolbar: false,
  }

  render() {
    const { children } = this.props;
    const toggleToolbar = () => {
      this.setState({ toolbar: !this.state.toolbar });
    };

    return (
      <MuiThemeProvider theme={muiTheme}>
        <div>
          <Toolbar toggleToolbar={toggleToolbar} />
          <Sidebar open={this.state.toolbar} />
          <div className={styles.content}>
            {children}
          </div>
        </div>
      </MuiThemeProvider>
    );
  }
}

export default createFragmentContainer(
  App,
  graphql`
    fragment App_viewer on Viewer {
      user {
        username
        hasPassword
      }
    }
  `,
);
