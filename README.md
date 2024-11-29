# frespect.com

# to copy
docker-compose up --build
 
## Before Production 

### 1. Remove volumes from docker-compose.yaml

### 2. Setup production mode for ReactJs app

### 3. Remove nodemon, it is not needed since we dont need live update

### 4. Configure skylla db for production

## To clone and develop project

### 1. Clone repository

### 2. in folder sakura-app create .env and paste:

```dotenv
REACT_APP_SOCKET_SERVER_HOST=localhost
REACT_APP_SOCKET_SERVER_PORT=5000
```

### 3. in folder sakura-media-server create .env and paste:

```dotenv
APP_PORT=4000
SIGNAL_SERVER_HOST=signal_server
SIGNAL_SERVER_PORT=5000
```

### 4. in folder sakura-signal create .env and paste: 

```dotenv
PORT=5000
MEDIA_SERVER_HOST=go_service
MEDIA_SERVER_API_KEY=your-secret-key
CORS_ORIGIN=*
REDIS_HOST=redis
REDIS_PORT=6379
```

### 5. Considering you've installed Docker Desktop run:

```bash
docker-compose up --build
```

### 6. You all set up!