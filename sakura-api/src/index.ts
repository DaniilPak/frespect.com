import 'reflect-metadata'; // Required for tsyringe
import dotenv from 'dotenv';
dotenv.config(); // Load environment variables

// Import the container file to ensure registration of dependencies
import { DownloadRoute } from './routes/download.route.js';

const PORT = process.env.PORT || 7000;

import cors from 'cors';
import express from 'express';
import { DownloadController } from './controllers/download.controller.js';
import { DownloadService } from './services/download.service.js';
import { DownloadedTrackRepository } from './repositories/downloaded-track.repository.js';
const app = express();

app.use(cors());
app.use(express.json());

import mongoose from 'mongoose';
import { MediaRoute } from './routes/media.route.js';
import { MediaController } from './controllers/media.controller.js';
import { MediaService } from './services/media.service.js';
const mongoHost = process.env.MONGO_DATABASE_HOST;

async function connectWithRetry() {
  try {
    console.log('MS', mongoHost);

    mongoose.connect(mongoHost!);
    const database = mongoose.connection;

    database.on('error', (error: any) => {
      console.log(error);
    });

    database.once('connected', () => {
      console.log('Database Connected');
    });
  } catch (err) {
    console.error('MongoDB connection failed, retrying in 5 seconds...', err);
    setTimeout(connectWithRetry, 5000); // Retry after 5 seconds
  }
}

connectWithRetry();

const downloadRoute = new DownloadRoute(
  new DownloadController(new DownloadService(new DownloadedTrackRepository()))
);

const mediaRoute = new MediaRoute(new MediaController(new MediaService()));

app.use('/download', downloadRoute.getRouter());
app.use('/media', mediaRoute.getRouter());

/// Main entry
app.listen(PORT, () => {
  console.log(`Server started at http://localhost:${PORT}`);
});
