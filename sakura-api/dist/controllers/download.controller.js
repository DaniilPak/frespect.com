var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import * as _ from 'kruza';
export class DownloadController {
    constructor(downloadService) {
        this.downloadService = downloadService;
        this.downloadService = downloadService;
    }
    downloadAudioFromYT(req, res) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const { ytlink } = req.body;
                if (!ytlink) {
                    res.status(400).send('Bad Request: Missing ytlink');
                    return;
                }
                // Proceed with your download logic using ytlink
                // ...
                yield this.downloadService.downloadAudio(ytlink);
                // download done
                _.log('Download complete');
                res.status(200).send('Download initiated');
            }
            catch (err) {
                console.error(err);
                res.status(500).send('Internal server error');
            }
        });
    }
}
