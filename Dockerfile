# Download go 1.14.4
FROM golang:1.14.4-alpine

# Update apt packages
RUN apk update

# Create jupyter workdir
ENV DIR /home/ubuntu/
WORKDIR $DIR

# Add ubuntu user (non-root)
RUN addgroup -S ubuntu && adduser -S ubuntu -G ubuntu

# Establish workdir as user directory
RUN chown -R ubuntu:ubuntu $DIR

# Switch to ubuntu user
USER ubuntu
