# Ctrlpad daemon

![GitHub stars](https://img.shields.io/github/stars/ctrlpad/daemon?style=for-the-badge&logo=github) ![GitHub issues](https://img.shields.io/github/issues/ctrlpad/daemon?style=for-the-badge&logo=github) ![Last commit](https://img.shields.io/github/last-commit/ctrlpad/daemon?style=for-the-badge&logo=github)

The daemon is the background service behind Ctrlpad. It scans for the `Ctrlpad_BLE` device, connects to it over Bluetooth Low Energy, subscribes to the button characteristic and executes the action that comes with every button press.

> [!WARNING]
> This project is still in development. Content, structure and setup can change at any time. The daemon currently supports Linux only - Windows and macOS are not supported at this time.

## Installation

Installation and service setup are documented in the guide:

**https://ctrlpad.xyz/docs/guide/ctrlpad-daemon**

It covers the install script, the Snap package, the Nix flake and the NixOS module, plus how to run the daemon as a service.

## Requirements

[devenv](https://devenv.sh) is required. It provides Go, BlueZ and GoReleaser in the versions this project is developed against. Install it together with Nix, then everything else comes from `devenv.nix`.

A working Bluetooth stack (BlueZ, running `bluetoothd`) is needed to talk to the device.

## Getting Started

```sh
git clone https://github.com/Ctrlpad/daemon.git
cd daemon

devenv shell
go run .
```

The daemon starts scanning immediately and logs every device it sees until it finds `Ctrlpad_BLE`.

For a watch-and-restart loop there is an [air](https://github.com/air-verse/air) config (`.air.toml`, builds to `./tmp/main`):

```sh
go install github.com/air-verse/air@latest
air
```

## Project Structure

```
main.go              entrypoint — logger setup, version info, calls cmd.Run
cmd/
  root.go            wires it together: connect → subscribe → execute in a loop
internal/
  ble/               BLE client
    client.go        scan, connect, enable notifications
    config.go        adapter, service and characteristic UUIDs
  executor/          turns a button payload into an action
    executor.go      button JSON + "type:target" action parsing
    linux.go         Linux action dispatch
  action/            the actions themselves
    application.go   launch a binary
    volume.go        set/raise/lower the volume via wpctl
  info/              version, commit and build date reporting
  utils/             small helpers (binary lookup in PATH)
nix/
  package.nix        buildGoModule package
  module.nix         NixOS module (systemd user service)
install.sh           install script for the released binaries
.goreleaser.yaml     release build: tar.gz archives + Snap
.github/workflows/   release workflow, runs GoReleaser on a published release
```

### How an action flows through it

The device sends a JSON button config over the notify characteristic:

```json
{ "id": 1, "name": "Firefox", "action": "application:firefox", "icon": "firefox" }
```

`executor` splits `action` into a type and a target, and dispatches on the type. Adding a new action means adding a function in `internal/action/` and a `case` in `internal/executor/linux.go`.

## Packaging

Releases are built by GoReleaser and published from `.github/workflows/release.yml` when a GitHub release is published. To check a build locally:

```sh
goreleaser release --snapshot --clean
```

The Nix flake exposes the `ctrlpad-daemon` package, an overlay and a NixOS module:

```sh
nix build .#ctrlpad-daemon
```

## Contributing

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/<you>/daemon.git`
3. **Branch**: `git checkout -b feature/your-feature`
4. **Commit**: `git commit -m 'feat: add some feature'`
5. **Push**: `git push origin feature/your-feature`
6. **Open** a PR

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org) — the changelog is generated from them (`feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`/`ci`; `chore` is filtered out).

Please follow the existing code style. Thanks!
