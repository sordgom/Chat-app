import React, { useEffect, useRef } from 'react';
import axios from 'axios';
import { useLocation } from 'react-router-dom';
import './VideoCall.css';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

const CallComponent = () => {
  const senderVideoRef = useRef(null);
  const receiverVideoRef = useRef(null);
  const pcSenderRef = useRef(null);
  const pcReceiverRef = useRef(null);
  let query = useQuery();

  useEffect(() => {
    const meetingId = query.get("meetingId");
    const peerId = query.get("peerId");
    const userId = query.get("userId");

    pcSenderRef.current = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    });

    pcReceiverRef.current = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    });

    pcSenderRef.current.onicecandidate = event => {
      if (event.candidate === null) {
        axios.post('/webrtc/sdp/m/' + meetingId + "/c/"+ userId + "/p/" + peerId + "/s/" + true,
        {"sdp" : btoa(JSON.stringify(pcSenderRef.current.localDescription))}).then(response => {
          pcSenderRef.current.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.data.Sdp))));
        });
      }
    };

    pcReceiverRef.current.onicecandidate = event => {
      if (event.candidate === null) {
        axios.post('/webrtc/sdp/m/' + meetingId + "/c/"+ userId + "/p/" + peerId + "/s/" + false, 
        {"sdp" : btoa(JSON.stringify(pcReceiverRef.current.localDescription))}).then(response => {
          pcReceiverRef.current.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.data.Sdp))));
        });
      }
    };
  }, []);

  const startCall = () => {
    navigator.mediaDevices.getUserMedia({video: true, audio: true}).then((stream) =>{
      senderVideoRef.current.srcObject = stream;
      var tracks = stream.getTracks();
      for (var i = 0; i < tracks.length; i++) {
        pcSenderRef.current.addTrack(stream.getTracks()[i]);
      }
      pcSenderRef.current.createOffer().then(d => pcSenderRef.current.setLocalDescription(d));
    });
    pcSenderRef.current.addEventListener('connectionstatechange', event => {
      if (pcSenderRef.current.connectionState === 'connected') {
        console.log("horray!")
      }
    });

    pcReceiverRef.current.addTransceiver('video', {'direction': 'recvonly'});
    pcReceiverRef.current.createOffer().then(d => pcReceiverRef.current.setLocalDescription(d));

    pcReceiverRef.current.ontrack = function (event) {
      receiverVideoRef.current.srcObject = event.streams[0];
      receiverVideoRef.current.autoplay = true;
      receiverVideoRef.current.controls = true;
    };
  };

  return (
    <div>
      <button onClick={startCall} className="start-call">Start the call!</button>
      <div className="container_row">
        <video autoPlay ref={senderVideoRef} width="500" height="500" controls muted></video>
        <div className="layer2">
          <video autoPlay ref={receiverVideoRef} width="160" height="120" controls muted></video>
        </div>
      </div>
    </div>
  );
};

export default CallComponent;
