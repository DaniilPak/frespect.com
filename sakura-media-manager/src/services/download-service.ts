import { Audio } from 'yt-converter';
import ffmpeg from 'fluent-ffmpeg';
import fs from 'fs';
import path from 'path';
const cassandra = require('cassandra-driver');

const SCYLLA_HOST = process.env.SCYLLA_HOST;

export class DownloadService {
  client: any;

  constructor() {
    // Initialize Cassandra client
    this.client = new cassandra.Client({
      contactPoints: [SCYLLA_HOST], // Replace with your ScyllaDB node IPs
      localDataCenter: 'datacenter1', // Replace with your data center name
      keyspace: 'your_keyspace', // Replace with your keyspace name
    });

    // Test the connection
    this.client
      .connect()
      .then(() => {
        console.log('Connected to ScyllaDB successfully.');
      })
      .catch((err: any) => {
        console.error('Failed to connect to ScyllaDB:', err);
      });

    const selectQuery =
      'SELECT id, name, value FROM test WHERE name = ? ALLOW FILTERING';
    const selectParams = ['Sample Name'];

    this.client
      .execute(selectQuery, selectParams, { prepare: true })
      .then((result: any) => {
        const row = result.first();
        if (row) {
          console.log('Retrieved row:', row);
          // Access individual fields if needed
          console.log('ID:', row.id);
          console.log('Name:', row.name);
          console.log('Value:', row.value);
        } else {
          console.log('No row found.');
        }
      })
      .catch((err: any) => {
        console.error('Failed to retrieve row:', err);
      });
  }

  async downloadAudio(url: string): Promise<void> {
    try {
      // Download the audio as an MP3 file
      const data = await Audio({
        url,
        onDownloading: (progress) => console.log(progress),
      });

      // Define the 'tracks' and 'tracks_ogg' directories
      // Move two directories up from __dirname
      const parentDir = path.resolve(__dirname, '../../');

      // Define the 'tracks_ogg' directory at the same level as 'tracks'
      const tracksOggDir = path.join(parentDir, 'tracks_ogg');

      // Define the input and output file paths
      const inputFilePath = data.pathfile;
      const outputFileName = `${path.basename(inputFilePath, '.mp3')}.ogg`;
      const outputFilePath = path.join(tracksOggDir, outputFileName);

      // Convert the MP3 file to OGG format
      ffmpeg(inputFilePath)
        .toFormat('ogg')
        .on('end', () => {
          console.log(`Conversion to OGG completed: ${outputFilePath}`);
        })
        .on('error', (err) => {
          console.error('Error during conversion:', err);
        })
        .save(outputFilePath);
    } catch (error) {
      console.error('Error downloading audio:', error);
    }
  }
}
