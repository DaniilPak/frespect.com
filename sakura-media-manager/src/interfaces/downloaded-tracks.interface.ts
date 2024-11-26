export interface DownloadedTrack {
  id: string; // UUID
  internal_id: string;
  video_id: string;
  title: string;
  file_path: string;
  duration: number;
  downloaded_at: Date;
}
