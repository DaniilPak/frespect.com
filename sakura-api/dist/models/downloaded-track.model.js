import mongoose from 'mongoose';
// Define the schema
const downloadedTrackSchema = new mongoose.Schema({
    sakura_audio_id: { type: String, required: true, unique: true }, // Primary Key
    video_id: { type: String, required: true }, // YouTube Video ID
    title: { type: String, required: true }, // Audio Title
    file_path: { type: String, required: true }, // File Path
    duration: { type: Number, required: true }, // Duration in seconds
}, { timestamps: true });
// Create the model
export const DownloadedTrackModel = mongoose.model('DownloadedTrack', downloadedTrackSchema);
