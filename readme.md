# Termfolio
> SSH-based interactive portfolio application served over SSH, built with Go, Wish, and Bubble Tea.

![20260209114926](https://cdn.tosh1ki.de/assets/images/20260209114926.png)
## 1: Project overview
### 1.1: What this project does
This project runs a terminal user interface over SSH so visitors can browse a personal portfolio without a browser.
The application includes menu-driven sections for about, projects, education, contact, privacy controls, and a live RSS feed view.

### 1.2: Main features
- SSH server using Wish and Bubble Tea.
- Keyboard-driven TUI with themed styling and animated logo.
- Privacy page that lets a visitor opt in or out of IP-based visit tracking.
- SQLite-backed unique visitor counter with opt-out persistence.
- RSS feed page that fetches and caches posts from `https://note.toshiki.dev/feed.xml`.

### 1.3: Navigation and keybinds
- `up` and `down` or `j` and `k`: move selection.
- `enter` or `space`: open selected page or confirm choice.
- `esc` or `backspace`: return to menu.
- `t`: cycle theme.
- `q` or `ctrl+c`: quit from menu.

## 2: Requirements
### 2.1: Runtime requirements
- Go `1.24.2` or compatible Go `1.24.x`.
- `ssh-keygen` available in `PATH` for host key generation.

### 2.2: Local networking expectations
- Local development is easiest on a non-privileged port such as `2222`.
- If you bind to port `22`, elevated privileges or a free privileged port may be required.

### 2.3: Bind to privileged ports on Linux
If you want to run on a privileged port such as `22` without running the app as root, grant the built binary bind capability:

```bash
sudo setcap 'cap_net_bind_service=+ep' ./bin/termfolio
```

## 3: Local development
### 3.1: Prepare configuration
Use the example file for local development:

```bash
cp config.yaml.example config.local.yaml
```

### 3.2: Generate host key
Generate a host key before first run:

```bash
make keys
```

If the configured key file is missing, the app can also prompt to generate one during startup.

### 3.3: Run the server
Run with an explicit local config:

```bash
go run . -c config.local.yaml
```

Or run the default command:

```bash
make run
```

### 3.4: Connect over SSH
Connect from another terminal:

```bash
ssh -p 2222 localhost
```

### 3.5: Common make targets
- `make run`: run application with `go run .`.
- `make build`: build binary to `bin/termfolio`.
- `make build-linux`: build Linux binary to `bin/termfolio-linux2`.
- `make keys`: create `.ssh/host_ed25519`.
- `make fmt`: run `go fmt ./...`.
- `make clean`: remove `bin/`.

## 4: Configuration
### 4.1: Config file structure
Default configuration keys:

```yaml
ssh:
  port: 2222
  address: "0.0.0.0"
  hostKeyPath: ".ssh/host_ed25519"

counter:
  enabled: true
  dbPath: "data/visitors.db"
```

The `counter` section supports either:
- Mapping form with `enabled` and optional `dbPath`.
- Scalar boolean form such as `counter: false`.

### 4.2: Environment variable overrides
These variables override file values:
- `SSH_PORT`
- `SSH_ADDRESS`
- `SSH_HOST_KEY_PATH`

### 4.3: Counter and privacy behavior
- Visitor count tracks unique IPs in SQLite.
- Opted-out IPs are stored in a dedicated table and removed from counted visitors.
- If tracking is disabled, the app still displays the current count without recording new visits.

## 5: Container and deployment
### 5.1: Docker image flow
- Multi-stage build compiles a static Linux binary.
- Runtime image is Alpine and includes `openssh-keygen`.
- Entry point generates host key when missing, then starts the app.

Build and run example:

```bash
docker build -t termfolio .
docker run --rm -p 2222:2222 termfolio
```

## 6: Repository structure
### 6.1: Key directories and files
```text
config/      configuration loading and defaults
counter/     SQLite visitor tracking store
pages/       TUI page renderers and content models
ui/          Bubble Tea app model and update loop
view/        theme palette and shared view helpers
main.go      SSH server bootstrap and middleware wiring
entrypoint.sh container startup and host key bootstrap
```

## 7: Version information
### 7.1: Current version constants
Version constants are defined in `version/version.go`.
At this snapshot:
- App name: `termfolio`
- Version: `0.1.4`

Print version:

```bash
go run . -v
```

## 8: License
This project is licensed under the MIT License.
See `/Users/andatoshiki/dev/github/termfolio/license` for the full text.
