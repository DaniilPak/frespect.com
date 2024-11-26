export {};

import { Router } from 'express';
import { DownloadController } from '../controllers/download-controller';
import { container } from 'tsyringe';

export class DownloadRoute {
  private readonly router: Router;
  private readonly downloadController: DownloadController;

  constructor() {
    this.router = Router();
    this.downloadController = container.resolve(DownloadController);
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
