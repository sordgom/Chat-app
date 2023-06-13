import React from "react";
import "./Header.scss";
import { Link } from 'react-router-dom';

const Header = () => (
  <div className="header">
    <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
      <h2>Realtime Chat App</h2>
    </Link>   
  </div>
);

export default Header;