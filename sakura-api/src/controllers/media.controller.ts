import { Request, Response } from 'express';
import { MediaService } from '../services/media.service.js';
import * as _ from 'kruza';

export class MediaController {
  constructor(private mediaService: MediaService) {
    this.mediaService = mediaService;
  }

  async downloadAudioFromYT(req: Request, res: Response): Promise<void> {
    try {
      const media_sid = req.params.sid;

      const mediaPath = await this.mediaService.getMedia(media_sid);

      res.sendFile(mediaPath);
    } catch (err) {
      console.error(err);
      res.status(500).send('Internal server error');
    }
  }
}
