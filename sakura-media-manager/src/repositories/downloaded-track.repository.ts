import { DownloadedTrack } from '../interfaces/downloaded-tracks.interface.js';
import { IRepository } from '../interfaces/IRepository.interface.js';

export class DownloadedTrackRepository implements IRepository<DownloadedTrack> {
  get(id: string): Promise<DownloadedTrack> {
    throw new Error('Method not implemented.');
  }
  getAll(): Promise<DownloadedTrack[]> {
    throw new Error('Method not implemented.');
  }
  create(item: DownloadedTrack): Promise<DownloadedTrack> {
    throw new Error('Method not implemented.');
  }
  update(id: string, item: DownloadedTrack): Promise<DownloadedTrack> {
    throw new Error('Method not implemented.');
  }
  delete(id: string): Promise<void> {
    throw new Error('Method not implemented.');
  }
}
