import path from 'path';
import { fileURLToPath } from 'url';
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
export class MediaService {
    constructor() { }
    getMedia(sakura_audio_id) {
        const mediaPath = path.join(__dirname, '../..', 'tracks_ogg', `${sakura_audio_id}.ogg`);
        return mediaPath;
    }
}
