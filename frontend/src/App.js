import { Box, ChakraProvider } from '@chakra-ui/react';
import React, { Component } from "react";
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import "./App.css";
import Chat from "./components/Chat/Chat";
import Footer from "./components/Footer/Footer";
import Header from './components/Header/Header';
import Home from "./components/Home/Home";
import Landing from "./components/Landing/Landing";
import Login from "./components/Login/Login";
import Register from "./components/Register/Register";
import VideoCallWrapper from "./components/VideoCall/VideoCallWrapper";
import theme from './theme';

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
              <Route path="/videochat" element={<VideoCallWrapper />} />
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
