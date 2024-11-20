// index.js
require("dotenv").config(); // Load environment variables

const express = require("express");
const http = require("http");
const socketIo = require("socket.io");
const axios = require("axios");
const bodyParser = require("body-parser");
const _ = require("kruza");

const app = express();

// Middleware to parse JSON bodies
app.use(bodyParser.json());

// Create HTTP server
const server = http.createServer(app);

// Initialize Socket.IO
const io = socketIo(server, {
  cors: {
    origin: process.env.CORS_ORIGIN || "*", // Set specific origin in production
    methods: ["GET", "POST"],
  },
});

// Configuration
const PORT = process.env.PORT || 5000;
const MEDIA_SERVER_HOST =
  process.env.MEDIA_SERVER_HOST || "http://localhost:4000";
const MEDIA_SERVER_API_KEY =
  process.env.MEDIA_SERVER_API_KEY || "your-secret-key";

// Mapping of clientId to Socket.IO socket
const clientSocketMap = new Map();

// Middleware to verify media server requests
const verifyMediaServer = (req, res, next) => {
  const apiKey = req.headers["x-api-key"];
  if (apiKey && apiKey === MEDIA_SERVER_API_KEY) {
    next();
  } else {
    res.status(401).json({ message: "Unauthorized" });
  }
};

// Socket.IO connection handler
io.on("connection", (socket) => {
  _.log(`Client connected: ${socket.id}`);

  clientSocketMap.set(socket.id, socket);
  _.log(`Registered clientId ${socket.id} with socket ${socket.id}`);

  // Handle SDP Offer from client
  socket.on("sdp-offer", async (data) => {
    try {
      // Forward the SDP offer to the media server
      await axios.post(
        `http://${MEDIA_SERVER_HOST}:4000/api/mediaserver`,
        {
          sdp: data.sdp,
          clientId: socket.id,
          roomId: "e7f31a47-3f2d-4bbf-8df4-9c2152f3b2a1"
          
        },
        {
          headers: {
            "x-api-key": MEDIA_SERVER_API_KEY, // If required by media server
          },
        }
      );
  
      // The media server is expected to send back a POST request to our microservice
    } catch (error) {
      console.error("Error processing SDP Offer:", error.message);
      socket.emit("error", { message: "Failed to process SDP Offer" });
    }
  });

  // Handle disconnection
  socket.on("disconnect", () => {
    _.log(`Client disconnected: ${socket.id}`);
    // Remove from clientSocketMap
    for (let [key, value] of clientSocketMap.entries()) {
      if (value.id === socket.id) {
        clientSocketMap.delete(key);
        break;
      }
    }
    // Notify media server about disconnection
    axios
      .post(
        `http://${MEDIA_SERVER_HOST}:4000/disconnect`,
        { clientId: socket.id },
        {
          headers: {
            "x-api-key": MEDIA_SERVER_API_KEY, // If required by media server
          },
        }
      )
      .catch((err) =>
        console.error("Error notifying media server:", err.message)
      );
  });
});

// --- Express POST Routes to Receive from Media Server ---

/**
 * Route: POST /media-server/answer
 * Description: Receive SDP answer from media server and send it to the appropriate client
 * Expected Body: { clientId: string, sdpAnswer: string }
 */
app.post("/media-server/answer", verifyMediaServer, (req, res) => {
  const { clientId, sdpAnswer } = req.body;

  if (!clientId || !sdpAnswer) {
    return res
      .status(400)
      .json({ message: "clientId and sdpAnswer are required" });
  }

  const clientSocket = clientSocketMap.get(clientId);

  if (clientSocket) {
    clientSocket.emit("sdp-answer", { sdp: sdpAnswer });
    console.log(`Sent SDP answer to client ${clientId}`);
    res.status(200).json({ message: "SDP answer sent to client" });
  } else {
    console.warn(`Socket not found for clientId ${clientId}`);
    res.status(404).json({ message: "Client not connected" });
  }
});

/**
 * Route: POST /media-server/some-other-route
 * Description: Additional routes as needed
 * Implement similar logic for other types of messages from media server
 */

// Start the server
server.listen(PORT, () => {
  _.log(`Socket.IO server running on port ${PORT}`);
  _.log("Server is listening");
});

// Connect to Redis
const redis = require("redis");

// Create a Redis client

const redisClient = redis.createClient({
  url: `redis://${process.env.REDIS_HOST || "127.0.0.1"}:${process.env.REDIS_PORT || "6379"}`, // Construct URL dynamically
});

// Connect to Redis
redisClient.connect().catch((err) => {
  console.error("Redis connection error:", err.message);
  _.log("Redis");
});

// Event listener for successful connection
redisClient.on("connect", () => {
  _.log("Connected to Redis");
});

(async () => {
  try {
    // Example map-like object
    const exampleMap = {
      key1: "value1",
      key2: "value2",
      key3: "value3",
    };

    const hashKey = "exampleMap";

    // Store the map in Redis as a hash
    // for (const [field, value] of Object.entries(exampleMap)) {
    //   await redisClient.hSet(hashKey, field, value);
    // }

    // Retrieve the map from Redis
    const retrievedMap = await redisClient.hGetAll(hashKey);

    console.log("Retrieved map from Redis:", retrievedMap);
  } catch (error) {
    console.error("Error working with Redis:", error.message);
  }
})();



