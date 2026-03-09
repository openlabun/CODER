FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    g++ \
    build-essential \
    python3 \
    coreutils \
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /runner
COPY runner.sh /runner/
# Fix line endings (CRLF -> LF) for Windows compatibility
RUN sed -i 's/\r$//' /runner/runner.sh
RUN chmod +x /runner/runner.sh
ENTRYPOINT ["/runner/runner.sh"]
