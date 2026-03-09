FROM eclipse-temurin:17-jdk-jammy

# Install Python and other utilities
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /runner
COPY runner.sh /runner/
RUN chmod +x /runner/runner.sh && sed -i 's/\r$//' /runner/runner.sh

ENTRYPOINT ["/runner/runner.sh"]
