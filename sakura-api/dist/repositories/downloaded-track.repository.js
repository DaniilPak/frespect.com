var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import { DownloadedTrackModel } from '../models/downloaded-track.model.js';
export class DownloadedTrackRepository {
    get(id) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const track = yield DownloadedTrackModel.findOne({ id }).exec();
                if (!track) {
                    throw new Error(`DownloadedTrack with id ${id} not found`);
                }
                return track;
            }
            catch (error) {
                console.error(`Error fetching DownloadedTrack with id ${id}:`, error);
                throw error;
            }
        });
    }
    getAll() {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                return yield DownloadedTrackModel.find().exec();
            }
            catch (error) {
                console.error(`Error fetching all DownloadedTracks:`, error);
                throw error;
            }
        });
    }
    create(item) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const newTrack = new DownloadedTrackModel(item);
                return yield newTrack.save();
            }
            catch (error) {
                console.error(`Error creating DownloadedTrack:`, error);
                throw error;
            }
        });
    }
    update(id, item) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const updatedTrack = yield DownloadedTrackModel.findOneAndUpdate({ id }, item, { new: true }).exec();
                if (!updatedTrack) {
                    throw new Error(`DownloadedTrack with id ${id} not found for update`);
                }
                return updatedTrack;
            }
            catch (error) {
                console.error(`Error updating DownloadedTrack with id ${id}:`, error);
                throw error;
            }
        });
    }
    delete(id) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const result = yield DownloadedTrackModel.deleteOne({ id }).exec();
                if (result.deletedCount === 0) {
                    throw new Error(`DownloadedTrack with id ${id} not found for deletion`);
                }
            }
            catch (error) {
                console.error(`Error deleting DownloadedTrack with id ${id}:`, error);
                throw error;
            }
        });
    }
}
