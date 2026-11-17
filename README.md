# WSJT-X Web Interface

A modern web interface for monitoring and controlling [WSJT-X](https://physics.princeton.edu/pulsar/k1jt/wsjtx.html) remotely.

**AI / agent contributors:** start at [`AGENTS.md`](./AGENTS.md) and work the queue in [`todo/`](./todo/).

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
    - Ensure `UDP Server` is set to `127.0.0.1:2237` (or the IP/port where `wsjtx-web` is listening).

### Using with GridTracker

[GridTracker](https://gridtracker.org/) (and GridTracker2) also listens for WSJT-X UDP on port **2237** by default. Only one process can bind that port, so if GridTracker is already running you will see:

```text
listen udp :2237: bind: address already in use
```

Pick one of these setups:

#### Option A — Run both (recommended): give wsjtx-web its own port

1. Leave GridTracker on its default `2237`.
2. Start wsjtx-web on a free port:

    ```bash
    go run main.go -udp :2238 -http :8080
    ```

3. Point WSJT-X at **wsjtx-web** for this session:
    - `File` > `Settings` > `Reporting`
    - Set `UDP Server` to `127.0.0.1:2238`
    - Keep `Accept UDP requests` enabled so reply/halt from the web UI still work

**Trade-off:** WSJT-X sends to one unicast destination. While it points at `:2238`, GridTracker on `:2237` will not receive that stream. Switch the WSJT-X UDP port back to `2237` when you want GridTracker only, or use Option B.

#### Option B — Run both via multicast

WSJT-X can send to a multicast group so multiple apps receive the same traffic:

1. In WSJT-X `Reporting`, set `UDP Server` to a multicast address such as `224.0.0.1` (port `2237`).
2. Configure GridTracker to join that same multicast group/port (see GridTracker’s UDP / multicast settings).
3. Start wsjtx-web joined to the same group:

    ```bash
    go run main.go -udp 224.0.0.1:2237 -http :8080
    ```

If multicast does not bind cleanly on your OS/firewall, fall back to Option A.

#### Option C — One app at a time on 2237

Stop GridTracker (or change its UDP listen port), then run the default:

```bash
go run main.go
```

wsjtx-web will listen on `:2237` and serve the UI on `http://localhost:8080`.

### Running Locally

1. **Start the Backend**:

    ```bash
    # First-time / fresh clone: create a placeholder so go:embed succeeds
    mkdir -p dist && echo '<!DOCTYPE html><title>wsjtx-web</title>' > dist/index.html

    go run main.go
    ```

    Defaults: WSJT-X UDP on `:2237`, HTTP on `http://localhost:8080`.

    Flags:

    ```bash
    go run main.go -udp :2238 -http :8080
    ```

    Audio capture for the waterfall is **off by default** (avoids malgo/PulseAudio noise). Enable it when you need the waterfall:

    ```bash
    go run main.go -audio
    ```

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
