var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export class MediaController {
    constructor(mediaService) {
        this.mediaService = mediaService;
        this.mediaService = mediaService;
    }
    downloadAudioFromYT(req, res) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const media_sid = req.params.sid;
                const mediaPath = yield this.mediaService.getMedia(media_sid);
                res.sendFile(mediaPath);
            }
            catch (err) {
                console.error(err);
                res.status(500).send('Internal server error');
            }
        });
    }
}
