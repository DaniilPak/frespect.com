import 'reflect-metadata'; // Required for tsyringe
import dotenv from 'dotenv';
dotenv.config(); // Load environment variables

// Import the container file to ensure registration of dependencies
import { DownloadRoute } from './routes/download-route.js';

const PORT = process.env.PORT || 7000;

import cors from 'cors';
import express from 'express';
import { DownloadController } from './controllers/download-controller.js';
import { DownloadService } from './services/download.service.js';
import { DownloadedTrackRepository } from './repositories/downloaded-track.repository.js';
const app = express();

app.use(cors());
app.use(express.json());

import mongoose from 'mongoose';
const mongoString = process.env.MONGO_DATABASE_URL;

mongoose.connect(mongoString!);
const database = mongoose.connection;

database.on('error', (error: any) => {
  console.log(error);
});

database.once('connected', () => {
  console.log('Database Connected');
});

const testRoute = new DownloadRoute(
  new DownloadController(new DownloadService(new DownloadedTrackRepository()))
);

app.use('/api/download', testRoute.getRouter());

/// Main entry
app.listen(PORT, () => {
  console.log(`Server started at http://localhost:${PORT}`);
});
