FROM debian:bullseye-slim

# Instala OpenJDK 17 y Python (para scripts del runner)
RUN apt-get update && apt-get install -y openjdk-17-jdk python3 && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /runner
COPY runner.sh /runner/
RUN chmod +x /runner/runner.sh && sed -i 's/\r$//' /runner/runner.sh

ENTRYPOINT ["/runner/runner.sh"]
