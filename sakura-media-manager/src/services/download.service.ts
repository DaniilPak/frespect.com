import { Audio } from 'yt-converter';
import ffmpeg from 'fluent-ffmpeg';
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
        url,
        onDownloading: (progress) => console.log(progress),
      });

      // Define the 'tracks' and 'tracks_ogg' directories
      // Move two directories up from __dirname
      const parentDir = path.resolve(__dirname, '../../');

      // Define the 'tracks_ogg' directory at the same level as 'tracks'
      const tracksOggDir = path.join(parentDir, 'tracks_ogg');

      const sakura_audio_id = nanoid(11);

      // Define the input and output file paths
      const inputFilePath = data.pathfile;
      const outputFileName = `${sakura_audio_id}.ogg`;
      const outputFilePath = path.join(tracksOggDir, outputFileName);

      // Convert the MP3 file to OGG format
      ffmpeg(inputFilePath)
        .toFormat('ogg')
        .on('end', async () => {
          console.log(`Conversion to OGG completed: ${outputFilePath}`);

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
