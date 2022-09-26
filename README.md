# stai-bouncer

[![License](https://img.shields.io/github/license/st3ffn/chia-bouncer)](/LICENSE)
[![GO](https://img.shields.io/github/go-mod/go-version/St3ffn/chia-bouncer)](https://golang.org/)

![STAI](https://www.station-i.de/wp-content/uploads/2021/12/stai_global_coin_logo_4c_4.0.png)

Tiny CLI tool to remove unwanted connections from your Stai Node based on the Geo IP Location (Country). 
The Tool is written in golang and acts more like a command line wrapper around `stai show...`
and the cli tool `geoiplookup`

## Getting started

### Pre-requisites

- Linux, MacOS (never tried it on Windows)
- `git` installed
- `go 1.16` installed
- `stai` is installed
- `geoiplookup` is installed (see below for installation instructions)

### Build from Source

Clone the repository

```shell
git clone https://github.com/h3ll0ofr1end/stai-bouncer.git
cd stai-bouncer
```

Build the binary

```shell
go build
```

### Install geoiplookup and update database

The tool `geoiplookup` is required to perform the Geo IP Location lookup.
Ubuntu 18 and Ubuntu 20 users can simply install it via:

```shell
sudo apt-get install geoip-bin
```

The installed package contains a database which is pretty old. 
To update it, download the newest from [here](https://dl.miyuru.lk/geoip/dbip/country/dbip4.dat.gz)
and unpack to `/usr/share/GeoIP`

```shell
wget https://dl.miyuru.lk/geoip/dbip/country/dbip4.dat.gz
gunzip dbip4.dat.gz
sudo mv dbip4.dat /usr/share/GeoIP/GeoIP.dat
```

You can check the age of the Geo IP Location database by running
```shell
geoiplookup localhost -v
   GeoIP Country Edition: DBIPLite-Country-CSV_20210401 converted to legacy DB with sherpya/geolite2legacy by miyuru.lk
   ...
```
Now you are ready to go.

### Usage

If your stai executable is located in `$HOME/stai-blockchain/venv/bin/stai`, you can simply run:
```bash
# assumes stai executable is located in $HOME/stai-blockchain/venv/bin/stai
# removes all connections from mars
> stai-bouncer mars
```
To specify a custom path to your stai executable use `--staiexec` or `-e`
```bash
# custom defined stai executable
# removes all connections from "elon on mars"
> stai-bouncer -e /home/steffen/stai-blockchain/venv/bin/stai elon on mars
```
You can also add another filter to remove all nodes which have a lower or equal down speed (in MiB) than specified. 
This will be independent of the location filter. It can be done via `--down-threshold` or `-d`.
```bash
# custom defined stai executable
# remove all connections with down speed lower or equal than 1.52 MiB
# removes all connections from "elon on mars"
> stai-bouncer -e /home/steffen/stai-blockchain/venv/bin/stai -d 1.52 elon on mars
```
Call with `--help` or `-h` to see the help page 
```bash
> stai-bouncer -h

NAME:
   stai-bouncer - remove unwanted connections from your Stai Node based on Geo IP Location.

USAGE:
   stai-bouncer [-e STAI-EXECUTABLE] [-d DOWN-THRESHOLD] LOCATION
   stai-bouncer -e /stai-blockchain/venv/bin/stai -d 0.2 mars

DESCRIPTION:
   Tool will lookup connections via 'stai show -c', get ip locations via geoiplookup and remove nodes from specified LOCATION via 'stai show -r'

GLOBAL OPTIONS:
   --stai-exec STAI-EXECUTABLE, -e STAI-EXECUTABLE     STAI-EXECUTABLE. normally located inside the bin folder of your venv directory (default: $HOME/stai-blockchain/venv/bin/stai)
   --down-threshold DOWN-THRESHOLD, -d DOWN-THRESHOLD  DOWN-THRESHOLD defines the additional filter for minimal down speed in MiB for filtering. (default: not active)
   --help, -h                                          show help (default: false)

COPYRIGHT:
   GNU GPLv3

```

### Verification

If you are on Linux or MacOS and you want the see the locations for the current connections you can use the following:

```shell
# go stai directory
cd stai-blockchain
# enable venv
. ./activate
# calls stai to show all connections, filter for FULL_NODE, print the ip, 
# print the IP and call geoiplookup command with ip, replace string "GeoIP Country Edition" with ""
stai show -c | grep FULL_NODE | awk '{print $2}' | xargs -I {} sh -c 'echo {} $(geoiplookup {})' | sed 's! GeoIP Country Edition!!'
```

### Integration

The script can easily be integrated with cron. Simply open the users crontab via `crontab -e` and add the following line.

```shell
# run stai-bouncer every 2 hours and remove all connections from mars
0 */2 * * * /PATH/TO/stai-bouncer mars
```

## Kind gestures

If you like the tool, and you are open for a kind gesture. Thanks in advance. 

- XCH Address: xch18s8r9v4kpwdx2y8jks5ma4g2rmff0h9dtr5nkc6zmnk5kj6v0faqer6k9v

