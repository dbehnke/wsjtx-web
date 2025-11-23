# WSJT-X Web Interface

A modern web interface for monitoring and controlling [WSJT-X](https://physics.princeton.edu/pulsar/k1jt/wsjtx.html) remotely.

## Features

- **Real-time Monitoring**: View decodes and band activity in real-time.
- **Remote Control**:
  - **Reply**: Double-click a decode to initiate a QSO.
  - **Halt Tx**: Stop transmission immediately.
- **Responsive Design**: Works on desktop and mobile devices.

## Architecture

- **Backend**: Go (Golang) server that bridges WSJT-X UDP traffic to WebSockets.
- **Frontend**: Vue.js 3 + TypeScript application for the user interface.

## Getting Started

### Prerequisites

- WSJT-X installed and running.
- Go 1.21+
- Node.js 18+

### Configuration

1. **WSJT-X Settings**:
    - Go to `File` > `Settings` > `Reporting`.
    - Check `Accept UDP requests`.
    - Ensure `UDP Server` is set to `127.0.0.1:2237` (or the IP where `wsjtx-web` is running).

### Running Locally

1. **Start the Backend**:

    ```bash
    go run main.go
    ```

    The server will start on `http://localhost:8080`.

2. **Start the Frontend** (for development):

    ```bash
    cd wsjtx-web-ui
    npm install
    npm run dev
    ```

    Access the UI at `http://localhost:5173`.

### Building for Production

1. **Build Frontend**:

    ```bash
    cd wsjtx-web-ui
    npm run build
    ```

2. **Build Backend**:

    ```bash
    cd ..
    go build
    ```

    This will create a `wsjtx-web` binary that serves the built frontend from the `dist` directory.

## License

MIT
