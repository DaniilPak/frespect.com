# Use the Node.js image
FROM node:18

# Set the working directory
WORKDIR /app

# Install dependencies
COPY package.json ./
RUN npm install

# Copy all app files
COPY . .

# Expose the React development server port
EXPOSE 3000

# Start the development server
CMD ["npm", "start"]
