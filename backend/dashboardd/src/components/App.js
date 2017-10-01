import React, { Component } from 'react';
import PropTypes from 'prop-types';

import logo from '../assets/logo.svg';

import './App.css';

class App extends Component {
  static propTypes = {
    children: PropTypes.element.isRequired,
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          {this.props.children}
        </p>
      </div>
    );
  }
}

export default App;
