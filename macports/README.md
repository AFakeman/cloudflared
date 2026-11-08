# Catalina MacPorts package

This local port builds the `catalina-compat` fork with the **MacPorts Go 1.24
toolchain**, installs `/opt/local/bin/cloudflared`, and supplies a system launchd
service. It conflicts with the standard `cloudflared` port.

The Portfile pins commit `f22d237ae55cdfd02ab8302c6262037828521311`, not a moving
branch. Local, uncommitted changes are not included in the package.

## Install on the target Mac

Install MacPorts for your macOS release and Apple's Command Line Tools first.
Update the ports tree so `go-1.24` supplies at least **1.24.13**, as required by
this branch's `go.mod`.

```sh
sudo /opt/local/bin/port selfupdate
cd /path/to/cloudflared/macports/net/cloudflared-catalina
/opt/local/bin/port lint
sudo /opt/local/bin/port install
/opt/local/bin/cloudflared --version
```

Installation downloads the pinned source archive, not your working tree. Go
modules are downloaded during the build and verified against `go.sum`; this is
not an offline build. `GOTOOLCHAIN=local` prevents automatic toolchain downloads.

## Configure and start

No credentials are bundled and the service is not started automatically.
The examples below assume the default MacPorts prefix `/opt/local`.

For a locally managed tunnel, put its credentials JSON at
`/opt/local/etc/cloudflared/TUNNEL-UUID.json` and create
`/opt/local/etc/cloudflared/config.yml`:

```yaml
tunnel: TUNNEL-UUID
credentials-file: /opt/local/etc/cloudflared/TUNNEL-UUID.json

ingress:
  - hostname: example.com
    service: http://127.0.0.1:80
  - service: http_status:404
```

Alternatively, for a remotely managed tunnel, store the tunnel token in
`/opt/local/etc/cloudflared/token` and use this configuration:

```yaml
token-file: /opt/local/etc/cloudflared/token
```

Protect whichever credentials you use:

```sh
sudo chown -R root:wheel /opt/local/etc/cloudflared
sudo chmod 700 /opt/local/etc/cloudflared
sudo chmod 600 /opt/local/etc/cloudflared/config.yml
# Also chmod 600 the credentials JSON or token file.
```

Then start and enable the service at boot:

```sh
sudo /opt/local/bin/port load cloudflared-catalina
sudo launchctl print system/org.macports.cloudflared-catalina
sudo tail -n 50 /opt/local/var/log/cloudflared/cloudflared.err.log
```

The LaunchDaemon runs as root. Do not also run `cloudflared service install`:
MacPorts owns this service. To stop it and disable startup:

```sh
sudo /opt/local/bin/port unload cloudflared-catalina
```

## Optional: register as a local ports repository

To install by name from any directory, run `portindex` inside `macports/`, then
add its **absolute** URL to `/opt/local/etc/macports/sources.conf`, before the
standard MacPorts source:

```text
file:///absolute/path/to/cloudflared/macports
```

Keep this checkout at that path and make it readable by the `macports` build
user. You can then run:

```sh
sudo /opt/local/bin/port install cloudflared-catalina
```

This also allows Ansible's `community.general.macports` module to install the
package by name. Ansible can manage credentials/configuration separately and
use `port load`/`port unload` for the service.

## Updating the snapshot

After pushing a new commit to `catalina-compat`:

1. Update the SHA in `go.setup`.
2. Download `https://github.com/AFakeman/cloudflared/archive/NEW-SHA.tar.gz` and
   recalculate its `rmd160`, `sha256`, and byte size for `checksums`.
3. Update the binary's `main.Version` suffix and increment `version` (or
   `revision` for packaging-only changes).
4. Run `port lint` and install/test on Catalina. Regenerate `PortIndex` if using
   the local repository method.

For an existing installation, unload before upgrading and load afterward:

```sh
sudo /opt/local/bin/port unload cloudflared-catalina
cd /path/to/cloudflared/macports/net/cloudflared-catalina
sudo /opt/local/bin/port upgrade
sudo /opt/local/bin/port load cloudflared-catalina
```

Tunnel configuration and credentials are administrator-owned files and are not
installed or overwritten by this package.
