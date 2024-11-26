import { Request, Response } from 'express';
import { DownloadService } from '../services/download-service';
import { inject, injectable } from 'tsyringe';
const _ = require('kruza');

@injectable()
export class DownloadController {
  constructor(
    @inject('DownloadService') private downloadService: DownloadService
  ) {
    this.downloadService = downloadService;
  }

  async downloadAudioFromYT(req: Request, res: Response): Promise<void> {
    try {
      const { ytlink } = req.body;

      if (!ytlink) {
        res.status(400).send('Bad Request: Missing ytlink');
        return;
      }

      // Proceed with your download logic using ytlink
      // ...
      await this.downloadService.downloadAudio(ytlink);

      // download done
      _.log('Download complete');

      res.status(200).send('Download initiated');
    } catch (err) {
      console.error(err);
      res.status(500).send('Internal server error');
    }
  }
}
