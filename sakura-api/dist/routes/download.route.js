import { Router } from 'express';
export class DownloadRoute {
    constructor(downloadController) {
        this.downloadController = downloadController;
        this.router = Router();
        this.downloadController = downloadController;
        this.setupRoutes();
    }
    setupRoutes() {
        this.router.post('/ytaudio', this.downloadController.downloadAudioFromYT.bind(this.downloadController));
    }
    getRouter() {
        return this.router;
    }
}
