import { exec } from 'child_process';
import { nanoid } from 'nanoid';
import path from 'path';
import { fileURLToPath } from 'url';
import { DownloadedTrackRepository } from '../repositories/downloaded-track.repository.js';
import ytdl from '@distube/ytdl-core';
import fs from 'fs';
import { pipeline } from 'stream/promises';
import ffmpegPath from 'ffmpeg-static';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export class DownloadService {
  constructor(private downloadedTrackRepository: DownloadedTrackRepository) {}

  async downloadAudio(url: string) {
    try {
      const sakura_audio_id = nanoid(11);
      const tracksDir = path.resolve(__dirname, '../../tracks');
      const tracksOggDir = path.resolve(__dirname, '../../tracks_ogg');
      const inputFilePath = path.join(tracksDir, `${sakura_audio_id}.mp3`);
      const outputFileName = `${sakura_audio_id}.ogg`;
      const outputFilePath = path.join(tracksOggDir, outputFileName);

      // Ensure directories exist
      if (!fs.existsSync(tracksDir))
        fs.mkdirSync(tracksDir, { recursive: true });
      if (!fs.existsSync(tracksOggDir))
        fs.mkdirSync(tracksOggDir, { recursive: true });

      const audioStream = ytdl(url, { filter: 'audioonly' });
      const writeStream = fs.createWriteStream(inputFilePath);

      // Await the completion of the piping process
      await pipeline(audioStream, writeStream);

      console.log('Download finished');

      // Convert to OGG using ffmpeg-static
      const command = `"${ffmpegPath}" -i "${inputFilePath}" -c:a libopus -page_duration 20000 -vn "${outputFilePath}"`;
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
          sakura_audio_id,
          video_id: ytdl.getURLVideoID(url),
          title: path.basename(inputFilePath),
          file_path: outputFilePath,
          duration: 0, // You might want to set the actual duration
          downloaded_at: new Date(),
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
