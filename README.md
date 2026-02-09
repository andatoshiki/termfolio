# Joe SSH Portfolio

A terminal portfolio you can SSH into. Built with Go, Wish, and Bubble Tea.

## Quick start (local)

1. Generate a host key (one time):

```bash
make keys
```

2. Run the server:

```bash
make run
```

3. Connect from another terminal:

```bash
ssh -p 2222 localhost
```

## Make targets

- `make run`: run locally with `go run .`
- `make build`: build a local binary at `bin/joe-ssh`
- `make build-linux`: build a Linux binary at `bin/joe-ssh-linux`
- `make keys`: generate `.ssh/host_ed25519`
- `make fmt`: run `go fmt ./...`
- `make clean`: remove the `bin/` directory

## Configuration (current defaults)

Right now the server listens on:

- Address: `0.0.0.0:2222`
- Host key path: `/data/host_ed25519`

You can override both with env vars:

- `SSH_ADDR`
- `SSH_HOST_KEY_PATH`

## Fly.io (cheapest setup)

This repo includes a minimal `fly.toml` that:

- Runs the app on internal port `2222`
- Exposes public port `22` so users can `ssh yourdomain.com`
- Mounts a volume at `/data` for a persistent host key
- Uses the smallest shared VM size

Create the volume once:

```bash
fly volumes create ssh_keys --size 1 --region <region>
```

## Notes

- This repo intentionally ignores `.ssh/` because it contains private keys.
- Built binaries are ignored; use `bin/` for artifacts.
