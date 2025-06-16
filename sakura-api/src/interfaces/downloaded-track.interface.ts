export interface DownloadedTrack {
  sakura_audio_id: string;
  video_id: string;
  title: string;
  file_path: string;
  duration: number;
  createdAt?: Date; // Timestamp for when the record was created
  updatedAt?: Date; // Timestamp for when the record was last updated
}
