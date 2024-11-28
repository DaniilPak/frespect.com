export {};

import { Router } from 'express';
import { MediaController } from '../controllers/media.controller.js';

export class MediaRoute {
  private readonly router: Router;

  constructor(private mediaController: MediaController) {
    this.router = Router();
    this.mediaController = mediaController;
    this.setupRoutes();
  }

  private setupRoutes() {
    this.router.get(
      '/:sid',
      this.mediaController.downloadAudioFromYT.bind(this.mediaController)
    );
  }

  public getRouter(): Router {
    return this.router;
  }
}
