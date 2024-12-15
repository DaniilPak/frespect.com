# frespect.com

## to copy

docker-compose up --build

## To set credentials for AWS

```bash
aws configure
```

## To get AWS identity

```bash
aws sts get-caller-identity
```

## To Login docker with ECS

```bash
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 863518443131.dkr.ecr.us-east-2.amazonaws.com
```

## Before Production

### 1. Remove volumes from docker-compose.yaml

### 2. Setup production mode for ReactJs app

### 3. Remove nodemon, it is not needed since we dont need live update

### 4. Configure mongodb change from db test to prod or dev

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

### 5. in folder sakura-media-manager create .env and paste:

```dotenv
PORT=7000
MONGO_DATABASE_HOST=mongodb://localhost:27017
```

### 6. Considering you've installed Docker Desktop run:

```bash
docker-compose up --build
```

### 7. Docker images and AWS ECR

## I created 3 repositories on ECR using those commands

```bash
aws ecr create-repository --repository-name node-api-repo
aws ecr create-repository --repository-name signal-server-repo
aws ecr create-repository --repository-name go-service-repo
aws ecr create-repository --repository-name redis-repo
aws ecr create-repository --repository-name mongo-repo
```

Results:

```bash
863518443131.dkr.ecr.us-east-2.amazonaws.com/node-api-repo
863518443131.dkr.ecr.us-east-2.amazonaws.com/signal-server-repo
863518443131.dkr.ecr.us-east-2.amazonaws.com/go-service-repo
863518443131.dkr.ecr.us-east-2.amazonaws.com/redis-repo
863518443131.dkr.ecr.us-east-2.amazonaws.com/mongo-repo
```

They are for the production build, we still use local images to test locally.

### 8. AWS IAM Permissions

AdministratorAccess

AdministratorAccess-Amplify

### 9. Secrets to save

Besides saving .env files, We have to save a key and secret key in file daniilp_accessKeys

everything located in folder frespect-secrets must be saved before moving project to another machine

### 10. Pushing local images to ECR

First i run:

```bash
docker-compose build
```

Second i tag all services:

```bash
docker tag frespectcom-go_service:latest 863518443131.dkr.ecr.us-east-2.amazonaws.com/go-service-repo:latest
docker tag frespectcom-node_api:latest 863518443131.dkr.ecr.us-east-2.amazonaws.com/node-api-repo:latest
docker tag frespectcom-signal_server:latest 863518443131.dkr.ecr.us-east-2.amazonaws.com/signal-server-repo:latest
docker tag redis:latest 863518443131.dkr.ecr.us-east-2.amazonaws.com/redis-repo
docker tag mongo:latest 863518443131.dkr.ecr.us-east-2.amazonaws.com/mongo-repo
```

Then i push em:

```bash
docker push 863518443131.dkr.ecr.us-east-2.amazonaws.com/go-service-repo:latest
docker push 863518443131.dkr.ecr.us-east-2.amazonaws.com/node-api-repo:latest
docker push 863518443131.dkr.ecr.us-east-2.amazonaws.com/signal-server-repo:latest
docker push 863518443131.dkr.ecr.us-east-2.amazonaws.com/redis-repo
docker push 863518443131.dkr.ecr.us-east-2.amazonaws.com/mongo-repo
```

### 11. EKS, Kubernetes and AWS

I installed eksctl using chocolatey

Now i create cluster

```bash
eksctl create cluster --name frespect-cluster --region us-east-2 --nodegroup-name frespect-nodegroup --node-type t3.micro --nodes 1
```

note: expensive 0.10$/hour , only for production
