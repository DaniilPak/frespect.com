import { DownloadedTrack } from '../interfaces/downloaded-track.interface.js';
import { IRepository } from '../interfaces/IRepository.interface.js';
import { DownloadedTrackModel } from '../models/downloaded-track.model.js';

export class DownloadedTrackRepository implements IRepository<DownloadedTrack> {
  async get(id: string): Promise<DownloadedTrack> {
    try {
      const track = await DownloadedTrackModel.findOne({ id }).exec();
      if (!track) {
        throw new Error(`DownloadedTrack with id ${id} not found`);
      }
      return track as DownloadedTrack;
    } catch (error) {
      console.error(`Error fetching DownloadedTrack with id ${id}:`, error);
      throw error;
    }
  }

  async getAll(): Promise<DownloadedTrack[]> {
    try {
      return await DownloadedTrackModel.find().exec();
    } catch (error) {
      console.error(`Error fetching all DownloadedTracks:`, error);
      throw error;
    }
  }

  async create(item: DownloadedTrack): Promise<DownloadedTrack> {
    try {
      const newTrack = new DownloadedTrackModel(item);
      return await newTrack.save();
    } catch (error) {
      console.error(`Error creating DownloadedTrack:`, error);
      throw error;
    }
  }

  async update(id: string, item: DownloadedTrack): Promise<DownloadedTrack> {
    try {
      const updatedTrack = await DownloadedTrackModel.findOneAndUpdate(
        { id },
        item,
        { new: true }
      ).exec();
      if (!updatedTrack) {
        throw new Error(`DownloadedTrack with id ${id} not found for update`);
      }
      return updatedTrack as DownloadedTrack;
    } catch (error) {
      console.error(`Error updating DownloadedTrack with id ${id}:`, error);
      throw error;
    }
  }

  async delete(id: string): Promise<void> {
    try {
      const result = await DownloadedTrackModel.deleteOne({ id }).exec();
      if (result.deletedCount === 0) {
        throw new Error(`DownloadedTrack with id ${id} not found for deletion`);
      }
    } catch (error) {
      console.error(`Error deleting DownloadedTrack with id ${id}:`, error);
      throw error;
    }
  }
}
