FROM golang:1.14-buster

# Set up dependencies
ENV PACKAGES curl make git
ENV PATH=/root/.cargo/bin:$PATH

# Set working directory for the build
WORKDIR /usr/local/app

# Install minimum necessary dependencies
RUN apt update && apt install -y $PACKAGES

# Install Rust and wasm32 dependencies
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Add source files
COPY . .

# Build client
RUN make build
