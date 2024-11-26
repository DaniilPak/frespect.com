import { container } from 'tsyringe';
import { DownloadService } from './services/download-service';

// Register services and repositories in the container
container.register('DownloadService', {
  useClass: DownloadService,
});
