export {};

import { Router } from 'express';
import { DownloadController } from '../controllers/download.controller.js';

export class DownloadRoute {
  private readonly router: Router;

  constructor(private downloadController: DownloadController) {
    this.router = Router();
    this.downloadController = downloadController;
    this.setupRoutes();
  }

  private setupRoutes() {
    this.router.post(
      '/ytaudio',
      this.downloadController.downloadAudioFromYT.bind(this.downloadController)
    );
  }

  public getRouter(): Router {
    return this.router;
  }
}
