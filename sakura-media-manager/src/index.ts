import 'reflect-metadata'; // Required for tsyringe
require('dotenv').config(); // Load environment variables

// Import the container file to ensure registration of dependencies
import './container';
import { DownloadRoute } from './routes/download-route';

const PORT = process.env.PORT || 7000;

const cors = require('cors');
const express = require('express');
const app = express();

app.use(cors());
app.use(express.json());

const testRoute = new DownloadRoute();

app.use('/api/download', testRoute.getRouter());

/// Main entry
app.listen(PORT, () => {
  console.log(`Server started at http://localhost:${PORT}`);
});
