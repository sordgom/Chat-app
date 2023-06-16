import React from 'react';
import axios from 'axios';
import './VideoCall.css';
import {
  Container,
  Flex,
  Textarea,
  Box,
  FormControl,
  FormErrorMessage,
  InputGroup,
  InputRightElement,
  Button,
  Input,
  VStack, 
  HStack
} from '@chakra-ui/react';
class CallComponent extends React.Component {
  constructor(props) {
    super(props);
    this.senderVideoRef = React.createRef();
    this.receiverVideoRef = React.createRef();
    this.pcSenderRef = null;
    this.pcReceiverRef = null;
    this.state = {
      socketConn: '',
      username: '',
      message: '',
      to: '',
      isInvalid: false,
      endpoint: 'http://localhost:8080',
      contact: '',
      contacts: [],
      renderContactList: [],
      chats: [],
      chatHistory: [],
      msgs: [],
    };
  }

  componentDidMount() {
    const params = new URLSearchParams(this.props.location.search);
    const meetingId = params.get("meetingId");
    const userId = params.get("userId");

    this.pcSenderRef = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    });

    this.pcReceiverRef = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    });

    this.pcSenderRef.onicecandidate = event => {
      if (event.candidate === null) {
        axios.post('/webrtc/sdp/m/' + meetingId + "/c/" + userId + "/p/" + this.state.username + "/s/" + true,
          { "sdp": btoa(JSON.stringify(this.pcSenderRef.localDescription)) }).then(response => {
            this.pcSenderRef.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.data.Sdp))));
          });
      }
    };

    this.pcReceiverRef.onicecandidate = event => {
      if (event.candidate === null) {
        axios.post('/webrtc/sdp/m/' + meetingId + "/c/" + userId + "/p/" + this.state.username + "/s/" + false,
          { "sdp": btoa(JSON.stringify(this.pcReceiverRef.localDescription)) }).then(response => {
            this.pcReceiverRef.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.data.Sdp))));
          });
      }
    };
  }
   // on change of input, set the value to the message state
   onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  addContact = async e => {
    e.preventDefault();
    try {
      const res = await axios.post(`${this.state.endpoint}/verify-contact`, {
        username: this.state.contact,
      });

      console.log(res.data);
      if (!res.data.status) {
        this.setState({ isInvalid: true });
      } else {
        // reset state on success
        this.setState({ username: this.state.contact, isInvalid: false }, () => {
          this.startCall();
        });
      }
      
    } catch (error) {
      console.error(error);
    }
  };

  startCall = () => {
    navigator.mediaDevices.getUserMedia({ video: true, audio: true }).then((stream) => {
      this.senderVideoRef.current.srcObject = stream;
      var tracks = stream.getTracks();
      for (var i = 0; i < tracks.length; i++) {
        this.pcSenderRef.addTrack(stream.getTracks()[i]);
      }
      this.pcSenderRef.createOffer().then(d => this.pcSenderRef.setLocalDescription(d));
    });

    this.pcSenderRef.addEventListener('connectionstatechange', event => {
      if (this.pcSenderRef.connectionState === 'connected') {
        console.log("horray!")
      }
    });

    this.pcReceiverRef.addTransceiver('video', { 'direction': 'recvonly' });
    this.pcReceiverRef.createOffer().then(d => this.pcReceiverRef.setLocalDescription(d));

    this.pcReceiverRef.ontrack = function (event) {
      this.receiverVideoRef.current.srcObject = event.streams[0];
      this.receiverVideoRef.current.autoplay = true;
      this.receiverVideoRef.current.controls = true;
    }.bind(this);
  };

  render() {
    return (
      <VStack spacing={4}>
        <Box>
            <FormControl isInvalid={this.state.isInvalid}>
              <InputGroup size="md">
                <Input
                  variant="flushed"
                  type="text"
                  placeholder="Add Receiver"
                  name="contact"
                  value={this.state.contact}
                  onChange={this.onChange}
                />
                <InputRightElement width="6rem">
                  <Button
                    colorScheme={'cyan'}
                    h="2rem"
                    size="lg"
                    variant="solid"
                    type="submit"
                    onClick={this.addContact}
                  >
                    Add
                  </Button>
                </InputRightElement>
              </InputGroup>
              {!this.state.isContactInvalid ? (
                ''
              ) : (
                <FormErrorMessage>contact does not exist</FormErrorMessage>
              )}
            </FormControl>
          </Box>
        <HStack spacing={10}>
          <Box bg={'teal'} borderRadius="md" p={1}>
            <video autoPlay ref={this.senderVideoRef} width="500" height="800" controls muted></video>
          </Box>
          <Box bg={'teal'} borderRadius="md" p={1}>
            <video autoPlay ref={this.receiverVideoRef} controls muted></video>
          </Box>
        </HStack>
      </VStack>
    );
  }
};

export default CallComponent;
