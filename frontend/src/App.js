import React, { Component } from "react";
import "./App.css";
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ChakraProvider, Box } from '@chakra-ui/react';
import theme from './theme';
import Header from './components/Header/Header';
import Footer from "./components/Footer/Footer";
import VideoCall from "./components/VideoCall/VideoCall";
import Home from "./components/Home/Home";
import Landing from "./components/Landing/Landing";
import Login from "./components/Login/Login";
import Register from "./components/Register/Register";
import Chat from "./components/Chat/Chat";

class App extends Component {

  render() {
    return (
      <ChakraProvider theme={theme}>
        <Box textAlign="center" fontSize="xl">
          <BrowserRouter>
            <Header></Header>
            <Routes>
              <Route path="/" element={<Landing />} />
              <Route path="/register" element={<Register />} />
              <Route path="/login" element={<Login />} />
              <Route path="/home" element={<Home />} />
              <Route path="/videochat" element={<VideoCall />} />
              <Route path="/chat"  element={<Chat />} />
            </Routes>
            <Footer></Footer>
          </BrowserRouter>
        </Box>
      </ChakraProvider>
    );
  }
}

export default App;
