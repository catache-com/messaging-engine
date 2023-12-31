# Use the official MongoDB image as the base image
FROM mongo:latest

# Create the MongoDB user and database
ENV MONGO_INITDB_ROOT_USERNAME=catache
ENV MONGO_INITDB_ROOT_PASSWORD=catache
ENV MONGO_INITDB_DATABASE=messaging-engine

# Copy the MongoDB configuration file
COPY mongod.conf /etc/mongod.conf

# Expose the default MongoDB port
EXPOSE 27017

# Set the data directory
ENV MONGO_DATA_DIR=/data/db/messaging-engine

# Create a volume for MongoDB data
VOLUME $MONGO_DATA_DIR

# Start the MongoDB service
CMD ["mongod", "-f", "/etc/mongod.conf"]
