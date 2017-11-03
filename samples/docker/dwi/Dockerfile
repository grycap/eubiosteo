FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install -y build-essential zlib1g-dev

# Install Python.
RUN \
  apt-get update && \
  apt-get install -y python python-dev python-pip python-virtualenv && \
  rm -rf /var/lib/apt/lists/*

COPY ./ /home

CMD ["python", "/home/main.py"]
