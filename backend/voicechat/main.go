package voicechat

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v2"
)

type Client struct {
	ID   string
	Peer *webrtc.PeerConnection
}

type Meeting struct {
	Clients map[string]*Client
}

func NewMeeting() *Meeting {
	return &Meeting{
		Clients: make(map[string]*Client),
	}
}

const (
	rtcpPLIInterval = time.Second * 3
)

var meetings = make(map[string]*Meeting)

type Sdp struct {
	Sdp string
}

func SetupVoiceChat() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	router := gin.Default()

	peerConnectionMap := make(map[string]chan *webrtc.Track)

	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	router.POST("/voice/sdp/m/:meetingId/c/:userId/p/:peerId/s/:isSender", func(c *gin.Context) {

		meetingID := c.Param("meetingId")
		isSender, _ := strconv.ParseBool(c.Param("isSender"))
		userID := c.Param("userId")
		peerID := c.Param("peerId")

		meeting, exists := meetings[meetingID]
		if !exists {
			meeting = NewMeeting()
			meetings[meetingID] = meeting
		}

		var session Sdp
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offer := webrtc.SessionDescription{}
		Decode(session.Sdp, &offer)

		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			log.Fatal(err)
		}

		meeting.Clients[userID] = &Client{ID: userID, Peer: peerConnection}

		if !isSender {
			receiveTrack(peerConnection, peerConnectionMap, peerID)
		} else {
			createTrack(peerConnection, peerConnectionMap, userID)
		}

		peerConnection.SetRemoteDescription(offer)

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, Sdp{Sdp: Encode(answer)})
	})

	router.Run(":8083")
}

func receiveTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	peerID string) {
	if _, ok := peerConnectionMap[peerID]; !ok {
		peerConnectionMap[peerID] = make(chan *webrtc.Track, 1)
	}
	localTrack := <-peerConnectionMap[peerID]
	peerConnection.AddTrack(localTrack)
}

func createTrack(peerConnection *webrtc.PeerConnection,
	peerConnectionMap map[string]chan *webrtc.Track,
	currentUserID string) {

	audioTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeOpus, 1234, "audio", "pion")
	if err != nil {
		log.Fatal(err)
	}
	peerConnection.AddTrack(audioTrack)

	go func() {
		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := audioTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatal(readErr)
			}

			// Process audio data here, e.g., apply audio effects or perform audio analysis

			// Send audio data to the remote peer
			

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err := audioTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatal(err)
			}

			// You can also use the rtpBuf to send the audio data over the network
			// or use audioTrack.Write() to write to the track if you have other sources of audio data
		}
	}()
}
