import { Router } from 'express';
export class MediaRoute {
    constructor(mediaController) {
        this.mediaController = mediaController;
        this.router = Router();
        this.mediaController = mediaController;
        this.setupRoutes();
    }
    setupRoutes() {
        this.router.get('/:sid', this.mediaController.downloadAudioFromYT.bind(this.mediaController));
    }
    getRouter() {
        return this.router;
    }
}
