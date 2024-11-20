// src/services/socket.js
import { io } from "socket.io-client";
import _ from "kruza";

// Get host and port from environment variables, with defaults
const SOCKET_SERVER_HOST =
  process.env.REACT_APP_SOCKET_SERVER_HOST || "localhost";
const SOCKET_SERVER_PORT = process.env.REACT_APP_SOCKET_SERVER_PORT || "5000";

// Construct the full Socket.IO server URL
const SOCKET_SERVER_URL = `http://${SOCKET_SERVER_HOST}:${SOCKET_SERVER_PORT}`;

const socket = io(SOCKET_SERVER_URL, {
  transports: ["websocket"],
});

socket.on("connect", () => {
  _.log("Connected to Socket.IO server:", socket.id);
});

socket.on("disconnect", () => {
  _.log("Disconnected from Socket.IO server");
});

export default socket;
