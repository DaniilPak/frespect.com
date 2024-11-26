# Use Node.js 18.20.4 base image
FROM node:18.20.4

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json first for caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Build ts
RUN npm run build

# Copy the entire project into the container
COPY . .

# Expose the app's port
EXPOSE 7000

# Start the application in development mode
CMD ["npm", "run", "dev"]
