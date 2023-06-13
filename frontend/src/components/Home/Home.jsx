import React from 'react';
import { Link } from 'react-router-dom';
import './Home.scss';

const HomePage = () => {
  return (
    <div className="home-page">
      <div className="component-window">
        <h2>Video Chat</h2>
        <Link to="/call">
          <div className="link-content">Click here to start a Video Call</div>
        </Link>      </div>
      <div className="component-window">
        <h2>Real Time Chat</h2>
        <Link to="/chat">
          <div className="link-content">Click here to start a Chat</div>
        </Link>
      </div>
    </div>
  );
};

export default HomePage;