// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

import React, { MutableRefObject, useEffect, useRef, useState } from "react";
import socket from "./services/socket";

const App: React.FC = () => {
  const rtcpPeerConnection = useRef<RTCPeerConnection | null>(null);
  const videoRef: MutableRefObject<HTMLVideoElement | null> = useRef(null);
  const [logs, setLogs] = useState<string[]>([]);
  const localSessionDescription = useRef<HTMLTextAreaElement | null>(null);
  const remoteSessionDescription = useRef<HTMLTextAreaElement | null>(null);
  const remoteVideosRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const config: RTCConfiguration = {
      iceServers: [
        {
          urls: "stun:stun.l.google.com:19302",
        },
      ],
    };

    rtcpPeerConnection.current = new RTCPeerConnection(config);

    const getMedia = async () => {
      try {
        const stream: MediaStream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        });

        stream.getTracks().forEach((track) => {
          track.enabled = true;
          rtcpPeerConnection.current?.addTrack(track, stream);
        });

        const offer = await rtcpPeerConnection.current?.createOffer();
        rtcpPeerConnection.current?.setLocalDescription(offer);
      } catch (err) {
        console.log("Error in getMedia: ", err);
      }
    };

    getMedia();

    rtcpPeerConnection.current.onconnectionstatechange = (e) => {
      const connection = rtcpPeerConnection.current;
      if (connection) {
        addLog(connection.iceConnectionState);
      }
    };

    rtcpPeerConnection.current.onicecandidate = (event) => {
      if (event.candidate === null) {
        const localSessionDescription = btoa(
          JSON.stringify(rtcpPeerConnection.current?.localDescription)
        );

        updateTextAreaValue(localSessionDescription);

        socket.emit("sdp-offer", { sdp: localSessionDescription });
        console.log("Sent SDP offer to server");
      }
    };

    rtcpPeerConnection.current.ontrack = function (event) {
      const videoElement = document.createElement("video");
      videoElement.autoplay = true;
      videoElement.controls = true;

      if (remoteVideosRef.current) {
        remoteVideosRef.current.appendChild(videoElement);
        // Use a method to assign the media stream to the video element
        videoElement.srcObject = event.streams[0];
      }
    };

    socket.on("sdp-answer", async (data) => {
      console.log("Received SDP answer from server:", data.sdp);

      updateRemoteSDP(data.sdp);
      startSession();
    });

    return () => {
      if (rtcpPeerConnection.current) {
        rtcpPeerConnection.current.close();
        rtcpPeerConnection.current = null;
      }

      socket.off("sdp-answer");
    };
  }, []);

  const addLog = (log: string) => {
    if (log) {
      setLogs((prevLogs) => [...prevLogs, log]);
    }
  };

  const updateTextAreaValue = (newValue: string) => {
    if (localSessionDescription.current) {
      localSessionDescription.current.value = newValue;
    }
  };

  const updateRemoteSDP = (newValue: string) => {
    if (remoteSessionDescription.current) {
      remoteSessionDescription.current.value = newValue;
    }
  }

  const copySDP = () => {
    if (localSessionDescription.current) {
      const textToCopy = localSessionDescription.current.value;
      navigator.clipboard
        .writeText(textToCopy)
        .then(() => {
          addLog("Local session description copied to clipboard");
        })
        .catch((err) => {
          console.error("Failed to copy text: ", err);
          addLog("Failed to copy local session description");
        });
    }
  };

  const startSession = () => {
    const remoteSDP = remoteSessionDescription.current?.value;

    if (!remoteSDP) {
      // Covers null, undefined, and empty string cases
      return alert("Session Description must not be empty");
    }

    try {
      // Decode the base64 encoded SDP
      const decodedSDP = atob(remoteSDP);

      // Parse the decoded SDP
      const parsedSDP = JSON.parse(decodedSDP);

      // Set the remote description using the parsed SDP
      rtcpPeerConnection.current?.setRemoteDescription(parsedSDP);
    } catch (err) {
      // More descriptive error handling
      alert(`Failed to set remote description: ${err}`);
    }
  };

  return (
    <div>
      <textarea
        id="localSessionDescription"
        readOnly // Use `readOnly` in JSX (not `readonly`)
        ref={localSessionDescription} // Attach the ref to the textarea
        style={{ width: "100%", height: "100px" }}
      ></textarea>
      <button onClick={copySDP}>Copy SDP</button>

      <textarea
        id="remoteSessionDescription"
        ref={remoteSessionDescription} // Attach the ref to the textarea
        style={{ width: "100%", height: "100px" }}
      ></textarea>
      <button onClick={startSession}>Start</button>

      <div id="remoteVideos" ref={remoteVideosRef}></div>

      <h2>Logs</h2>
      <div
        id="logs"
        style={{
          width: "100%",
          height: "200px",
          border: "1px solid #ccc",
          overflowY: "scroll",
        }}
      >
        {logs.map((log, index) => {
          return <div key={index}>{log}</div>;
        })}
      </div>
    </div>
  );
};

export default App;
