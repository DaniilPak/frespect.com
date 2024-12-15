# Use Node.js 18.20.4 base image
FROM node:18.20.4

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json first for caching
COPY package.json package-lock.json ./

# For Dev
# RUN npm install

# Install dependencies
RUN npm ci --only=production

# Copy the entire project into the container
COPY . .

# Expose the app's port
EXPOSE 5000

# For Dev
# CMD ["npx", "nodemon", "index.js"]

# Start the application
CMD ["node", "index.js"]
