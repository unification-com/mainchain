# Running `und` as a background service

#### Contents

[[toc]]

If you intend to run your node as a Validator on any of the public networks, then you will most likely need to permanently run `und` as a background service (as opposed to manually running `und start` and leaving a terminal window/SSH session open).

This can easily be done using `systemctl`, and setting up an appropriate service configuration.

The following is a generic \*nix guide, and may need adapting for your particular distribution.

Any text editor can be used to create the service configuration file, for example `nano`:

```bash
sudo nano /etc/systemd/system/und.service
```

At a minimum, the service configuration should contain:

```
[Unit]
Description=Unification Mainchain Validator Node

[Service]
User=USERNAME
Group=USERNAME
WorkingDirectory=/home/USERNAME
ExecStart=/home/USERNAME/go/bin/und start --home /home/USERNAME/.und_mainchain
LimitNOFILE=4096

[Install]
WantedBy=default.target
```

Of course, `und` can also be installed into `/usr/local/bin` instead of `/home/USERNAME/$GOPATH/bin/und`

It is entirely possible to create a more sophisticated service definition should you desire.

Next, inform `systemctl` of the new service:

```bash
sudo systemctl daemon-reload
```

The service can now be started:

```bash
sudo systemctl start und
```

and stopped:

```bash
sudo systemctl stop und
```

in the background.

Finally, you can monitor the log output for the service by running:

```bash
$ sudo journalctl -u und --follow
```
