import React, { Component } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from './components/Header/Header';
import ChatHistory from "./components/ChatHistory";
import ChatInput from "./components/ChatInput/ChatInput";
import Footer from "./components/Footer/Footer";
import VideoCall from "./components/VideoCall/VideoCall";
import Home from "./components/Home/Home";
class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }

  send(event) {
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }

  render() {
    return (
      <div className="App">
        <Router>
        <Header />
          <Routes>
            <Route path="/" element={<Home/>} />
            <Route path="/chat" element={
              <div className="ChatHistory">
                <ChatHistory chatHistory={this.state.chatHistory} />
                <ChatInput send={this.send} />
              </div>
            } />
            <Route path="/call" element={<VideoCall />} />
          </Routes>
        <Footer />
      </Router>
      </div>
      
    );
  }
}

export default App;
