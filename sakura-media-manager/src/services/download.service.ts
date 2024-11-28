import { exec } from 'child_process';
import { Audio } from 'yt-converter';
import { nanoid } from 'nanoid';
import path from 'path';
import { fileURLToPath } from 'url';
import { DownloadedTrackRepository } from '../repositories/downloaded-track.repository.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export class DownloadService {
  constructor(private downloadedTrackRepository: DownloadedTrackRepository) {}

  async downloadAudio(url: string): Promise<void> {
    try {
      // Download the audio as an MP3 file
      const data = await Audio({
        directory: './tracks',
        url,
        onDownloading: (progress) => console.log(progress),
      });

      // Define the 'tracks_ogg' directory
      const parentDir = path.resolve(__dirname, '../../');
      const tracksOggDir = path.join(parentDir, 'tracks_ogg');

      // Generate a unique ID for the audio file
      const sakura_audio_id = nanoid(11);

      // Define the input and output file paths
      const inputFilePath = data.pathfile;
      const outputFileName = `${sakura_audio_id}.ogg`;
      const outputFilePath = path.join(tracksOggDir, outputFileName);

      // Run the ffmpeg command directly to convert to Opus format
      const command = `ffmpeg -i "${inputFilePath}" -c:a libopus -page_duration 20000 -vn "${outputFilePath}"`;

      exec(command, async (error, stdout, stderr) => {
        if (error) {
          console.error(`Error during conversion: ${error.message}`);
          return;
        }
        if (stderr) {
          console.error(`FFmpeg stderr: ${stderr}`);
        }

        console.log('Conversion finished');

        // Create a new record in the repository
        const trackRecord = {
          sakura_audio_id: sakura_audio_id, // Generate unique ID
          video_id: 'yt_id', // YouTube video ID
          title: path.basename(inputFilePath), // Audio title
          file_path: outputFilePath, // Path to OGG file
          duration: 0, // Audio duration in seconds
          downloaded_at: new Date(), // Current timestamp
        };

        try {
          await this.downloadedTrackRepository.create(trackRecord);
          console.log('Track record saved successfully:', trackRecord);
        } catch (err) {
          console.error('Error saving track record:', err);
        }
      });
    } catch (error) {
      console.error('Error downloading audio:', error);
    }
  }
}
