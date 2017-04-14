# Luzifer / expose

`expose` is a small CLI utility to control a running [ngrok 2](https://ngrok.com/) daemon.

## Use cases

- Quickly open a local port to demonstrate something to someone
    ```bash
    # expose create 4000
    Created tunnel to address "localhost:4000" with URL https://8d35e4bf.eu.ngrok.io
    ```
- Close that port again
    ```bash
    # go run *.go delete 4000
    Successfully closed tunnel "expose_4000" to address "localhost:4000".
    ```
- Quickly share a folder with files through HTTP
    ```bash
    # go run *.go serve
    Created HTTP server for this directory with URL https://81e668af.eu.ngrok.io
    Press Ctrl+C to stop server and tunnel
    ```
- List the active tunnels in a nice table
    ```bash
    # go run *.go list
         NAME     | TYPE  |    ADDRESS     |          PUBLIC URL
    +-------------+-------+----------------+------------------------------+
      expose_3000 | https | localhost:3000 | https://9c8d36d6.eu.ngrok.io
    ```

## How to set up

- Download the latest release of [ngrok 2](https://ngrok.com/)
- Configure ngrok (see the official documentation for that!)
  - For this to work you only need an `authtoken` set in `~/.ngrok2/ngrok.yml`
- Set it up to start automatically (see `docs/ngrok.service` for a systemd service file)
- Download the [latest release](https://github.com/Luzifer/expose/releases/latest) of this tool
