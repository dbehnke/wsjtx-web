# Project Instructions

## Overview

WSJT-X Web Interface is a remote control and monitoring interface for WSJT-X. It consists of a Go backend that communicates with WSJT-X via UDP and a Vue.js frontend for the user interface.

## Architecture

- **Backend**: Go (Golang). Handles UDP communication with WSJT-X and serves the frontend via HTTP/WebSocket.
- **Frontend**: Vue.js 3 + TypeScript + Vite. Connects to the backend via WebSocket to receive status updates and decodes, and to send commands.

## Development

- **Backend**:
  - Run: `go run main.go`
  - Build: `go build`
- **Frontend**:
  - Directory: `wsjtx-web-ui`
  - Install: `npm install`
  - Dev: `npm run dev`
  - Build: `npm run build`

## Testing

- Run backend tests: `go test ./...`
- Run frontend tests: `npm run test:unit`

## Deployment

- Build frontend: `cd wsjtx-web-ui && npm run build`
- Build backend: `go build`
- The backend embeds the `dist` folder from the frontend build.
