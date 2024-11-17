# Use the Node.js 18.20.4 base image
FROM node:18.20.4

# Set working directory
WORKDIR /app

# Copy package.json and install dependencies
COPY package.json ./
RUN npm install

# Copy project files
COPY . .

# Expose the app's port
EXPOSE 5000

# Start the app
CMD ["node", "index.js"]
