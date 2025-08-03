import React, { useState, useEffect, useRef } from "react";
import {
  Box,
  Button,
  TextField,
  Typography,
  Card,
  CardContent,
  Grid,
  Alert,
  CircularProgress,
  IconButton,
  Paper,
} from "@mui/material";
import { get } from "../api/fetch";
import { Videocam, VideocamOff, Mic, MicOff } from "@mui/icons-material";

const WebRTCClient = () => {
  const [localStream, setLocalStream] = useState(null);
  const [remoteStream, setRemoteStream] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);
  const [error, setError] = useState(null);
  const [roomId, setRoomId] = useState("");
  const [isVideoEnabled, setIsVideoEnabled] = useState(true);
  const [isAudioEnabled, setIsAudioEnabled] = useState(true);
  const [connectionState, setConnectionState] = useState("disconnected");

  const localVideoRef = useRef(null);
  const remoteVideoRef = useRef(null);
  const peerConnectionRef = useRef(null);
  const signalingSocketRef = useRef(null);
  const localStreamRef = useRef(null);

  // WebRTC configuration
  const rtcConfiguration = {
    iceServers: [
      { urls: "stun:stun.l.google.com:19302" },
      { urls: "stun:stun1.l.google.com:19302" },
      { urls: "stun:stun2.l.google.com:19302" },
      { urls: "stun:stun3.l.google.com:19302" },
      { urls: "stun:stun4.l.google.com:19302" },
    ],
  };

  // Initialize local media stream
  const initializeLocalStream = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: isVideoEnabled,
        audio: isAudioEnabled,
      });

      localStreamRef.current = stream;
      setLocalStream(stream);

      if (localVideoRef.current) {
        localVideoRef.current.srcObject = stream;
      }

      setError(null);
    } catch (err) {
      setError(`Failed to access media devices: ${err.message}`);
      console.error("Error accessing media devices:", err);
    }
  };

  // Create peer connection
  const createPeerConnection = () => {
    const pc = new RTCPeerConnection(rtcConfiguration);

    // Add local stream tracks to peer connection
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach((track) => {
        pc.addTrack(track, localStreamRef.current);
      });
    }

    // Handle incoming remote stream
    pc.ontrack = (event) => {
      console.log("Received remote stream");
      setRemoteStream(event.streams[0]);
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = event.streams[0];
      }
    };

    // Handle connection state changes
    pc.onconnectionstatechange = () => {
      setConnectionState(pc.connectionState);
      console.log("Connection state:", pc.connectionState);

      if (pc.connectionState === "connected") {
        setIsConnected(true);
        setIsConnecting(false);
      } else if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        setIsConnected(false);
        setIsConnecting(false);
        setRemoteStream(null);
      }
    };

    // Handle ICE candidate events
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        sendSignalingMessage({
          type: "Candidate",
          data: event.candidate,
          // roomId: roomId
        });
        console.log("Candidate", event.candidate);
      }
    };

    peerConnectionRef.current = pc;
    return pc;
  };

  // Connect to signaling server
  const connectToSignalingServer = async () => {
    try {
      const wsURl = await get("/api/interviews");
      const socket = new WebSocket(wsURl.data);

      socket.onopen = () => {
        console.log("Connected to signaling server");
        setError(null);
      };

      socket.onmessage = async (event) => {
        const message = JSON.parse(event.data);
        await handleSignalingMessage(message);
      };

      socket.onerror = (error) => {
        setError(`Signaling server error: ${error.message}`);
        console.error("Signaling server error:", error);
      };

      socket.onclose = () => {
        console.log("Disconnected from signaling server");
        setIsConnected(false);
        setIsConnecting(false);
      };

      signalingSocketRef.current = socket;
    } catch (err) {
      setError(`Failed to connect to signaling server: ${err.message}`);
      console.error("Error connecting to signaling server:", err);
    }
  };

  // // Send message to signaling server
  const sendSignalingMessage = (message) => {
    if (
      signalingSocketRef.current &&
      signalingSocketRef.current.readyState === WebSocket.OPEN
    ) {
      signalingSocketRef.current.send(JSON.stringify(message));
    }
  };

  // Handle incoming signaling messages
  const handleSignalingMessage = async (message) => {
    const pc = peerConnectionRef.current;

    switch (message.type) {
      case "Offer":
        console.log("Received offer");
        await pc.setRemoteDescription(new RTCSessionDescription(message.offer));
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        sendSignalingMessage({
          type: "Answer",
          data: answer,
        });
        break;

      case "Answer":
        console.log("Received answer");
        await pc.setRemoteDescription(
          new RTCSessionDescription(message.answer)
        );
        break;

      case "Candidate":
        console.log("Received ICE candidate");
        if (pc.remoteDescription) {
          await pc.addIceCandidate(new RTCIceCandidate(message.candidate));
        }
        break;
      
      //   case 'user-joined':
      //     console.log('User joined room');
      //     if (pc.signalingState === 'stable') {
      //   const offer = await pc.createOffer();
      //       await pc.setLocalDescription(offer);
      //       sendSignalingMessage({
      //         type: 'offer',
      //         offer: offer,
      //         roomId: roomId
      //       });
      //     }
      //     break;

      default:
        console.log("Unknown message type:", message.type);
    }
  };

  // Join room
  const joinRoom = async () => {
    if (!roomId.trim()) {
      setError("Please enter a room ID");
      return;
    }

    setIsConnecting(true);
    setError(null);

    try {
      // Initialize local stream if not already done
      if (!localStreamRef.current) {
        await initializeLocalStream();
      }

      // Connect to signaling server
      await connectToSignalingServer();

      // Wait for socket to be ready
      setTimeout(() => {
        if (
          signalingSocketRef.current &&
          signalingSocketRef.current.readyState === WebSocket.OPEN
        ) {
          // Create peer connection
          createPeerConnection();

          // Join room
          // sendSignalingMessage({
          //   type: 'join-room',
          //   roomId: roomId
          // });
        } else {
          setError("Failed to connect to signaling server");
          setIsConnecting(false);
        }
      }, 1000);
    } catch (err) {
      setError(`Failed to join room: ${err.message}`);
      setIsConnecting(false);
      console.error("Error joining room:", err);
    }
  };

  // Leave room
  const leaveRoom = () => {
    if (signalingSocketRef.current) {
      signalingSocketRef.current.close();
    }

    if (peerConnectionRef.current) {
      peerConnectionRef.current.close();
    }

    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach((track) => track.stop());
      localStreamRef.current = null;
    }

    setLocalStream(null);
    setRemoteStream(null);
    setIsConnected(false);
    setIsConnecting(false);
    setConnectionState("disconnected");
    setError(null);
  };

  // Toggle video
  const toggleVideo = () => {
    if (localStreamRef.current) {
      const videoTrack = localStreamRef.current.getVideoTracks()[0];
      if (videoTrack) {
        videoTrack.enabled = !videoTrack.enabled;
        setIsVideoEnabled(videoTrack.enabled);
      }
    }
  };

  // Toggle audio
  const toggleAudio = () => {
    if (localStreamRef.current) {
      const audioTrack = localStreamRef.current.getAudioTracks()[0];
      if (audioTrack) {
        audioTrack.enabled = !audioTrack.enabled;
        setIsAudioEnabled(audioTrack.enabled);
      }
    }
  };

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      leaveRoom();
    };
  }, []);

  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: "auto" }}>
      <Typography variant="h4" gutterBottom>
        WebRTC Remote Client
      </Typography>

      {/* Connection Settings */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Connection Settings
          </Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Signaling Server"
                // value={signalingServer}
                // onChange={(e) => setSignalingServer(e.target.value)}
                disabled={isConnecting || isConnected}
                margin="normal"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Room ID"
                value={roomId}
                onChange={(e) => setRoomId(e.target.value)}
                disabled={isConnecting || isConnected}
                margin="normal"
              />
            </Grid>
          </Grid>
          <Box sx={{ mt: 2 }}>
            {!isConnected && !isConnecting ? (
              <Button
                variant="contained"
                color="primary"
                onClick={joinRoom}
                disabled={!roomId.trim()}
              >
                Join Room
              </Button>
            ) : (
              <Button variant="contained" color="error" onClick={leaveRoom}>
                Leave Room
              </Button>
            )}
          </Box>
        </CardContent>
      </Card>

      {/* Status and Error Display */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {isConnecting && (
        <Alert severity="info" sx={{ mb: 3 }}>
          <CircularProgress size={20} sx={{ mr: 1 }} />
          Connecting to room...
        </Alert>
      )}

      {isConnected && (
        <Alert severity="success" sx={{ mb: 3 }}>
          Connected to room: {roomId}
        </Alert>
      )}

      {/* Video Grid */}
      <Grid container spacing={3}>
        {/* Local Video */}
        <Grid item xs={12} md={6}>
          <Paper elevation={3} sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Local Video
            </Typography>
            <Box
              sx={{
                position: "relative",
                width: "100%",
                height: 300,
                backgroundColor: "#000",
                borderRadius: 1,
                overflow: "hidden",
              }}
            >
              <video
                ref={localVideoRef}
                autoPlay
                muted
                playsInline
                style={{
                  width: "100%",
                  height: "100%",
                  objectFit: "cover",
                }}
              />
              {!isVideoEnabled && (
                <Box
                  sx={{
                    position: "absolute",
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    backgroundColor: "rgba(0,0,0,0.7)",
                    color: "white",
                  }}
                >
                  <VideocamOff sx={{ fontSize: 60 }} />
                </Box>
              )}
            </Box>
          </Paper>
        </Grid>

        {/* Remote Video */}
        <Grid item xs={12} md={6}>
          <Paper elevation={3} sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Remote Video
            </Typography>
            <Box
              sx={{
                position: "relative",
                width: "100%",
                height: 300,
                backgroundColor: "#000",
                borderRadius: 1,
                overflow: "hidden",
              }}
            >
              {remoteStream ? (
                <video
                  ref={remoteVideoRef}
                  autoPlay
                  playsInline
                  style={{
                    width: "100%",
                    height: "100%",
                    objectFit: "cover",
                  }}
                />
              ) : (
                <Box
                  sx={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    height: "100%",
                    color: "white",
                  }}
                >
                  <Typography variant="body1">
                    Waiting for remote stream...
                  </Typography>
                </Box>
              )}
            </Box>
          </Paper>
        </Grid>
      </Grid>

      {/* Controls */}
      {isConnected && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Controls
            </Typography>
            <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
              <IconButton
                onClick={toggleVideo}
                color={isVideoEnabled ? "primary" : "error"}
                size="large"
              >
                {isVideoEnabled ? <Videocam /> : <VideocamOff />}
              </IconButton>
              <IconButton
                onClick={toggleAudio}
                color={isAudioEnabled ? "primary" : "error"}
                size="large"
              >
                {isAudioEnabled ? <Mic /> : <MicOff />}
              </IconButton>
            </Box>
          </CardContent>
        </Card>
      )}

      {/* Connection State */}
      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Connection Status
          </Typography>
          <Typography variant="body2" color="text.secondary">
            State: {connectionState}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Local Stream: {localStream ? "Active" : "Inactive"}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Remote Stream: {remoteStream ? "Active" : "Inactive"}
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
};

export default WebRTCClient;
