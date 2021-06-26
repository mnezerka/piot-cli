# PIOT Command Line Client

Command line client for managing PIOT infrastructure as well as for gathering
various data sets (e.g. export of sensor readings) 

# Installation

Download latest version of `piot` binary for your architecture. The recommended location
for saving is directory, which is registered in your search path (e.g. `~/bin`).

# Configuration

Tool needs several configuraion parameters to be able to access piot server or influx database:

| Parameter           | Environment variable     | Description                                                   |
|---------------------|--------------------------|---------------------------------------------------------------|
| `piot.url`          | `PIOT_PIOT_URL`          | URL of the PIOT server                                        |
| `piot.user`         | `PIOT_PIOT_USER`         | User for PIOT server                                          |
| `piot.password`     | `PIOT_PIOT_PASSWORD`     | Password for PIOT server                                      |
| `log.level`         | `PIOT_LOG_LEVEL`         | Log level for command line tool (DEBUG, INFO, WARNING, ERROR) |
| `influxdb.url`      | `PIOT_INFLUXDB_URL`      | URL of the Influx Database                                    |
| `influxdb.user`     | `PIOT_INFLUXDB_USER`     | User for Influx Database                                      |
| `influxdb.password` | `PIOT_INFLUXDB_PASSWORD` | Password for Influx Database                                  |

## Config file

The recommended way to provide all configuration parameters is yaml
configuration file. Default name is `.piot.yaml` and default search path is home
directory. You can override default location by specifying value for `--config`
flag. This is example of configuration file:

```
---
piot.url: https://example.com/api
piot.user: piotuser@example.com
piot.password: piotpassword
log.level: INFO
influxdb.url: https://example.com/influxdb
influxdb.user: influxuser
influxdb.password: influxpassword
```

## Environment variables

All parameters could be set also in shell environment variables. This is how to
change log level to `DEBUG` in linux shell:

```
export PIOT_LOG_LEVEL=DEBUG
```

Parameter values set in environment variables override values from configuration
file.

## Command line flags

Another way to specify parameter values is directly via command line flags
(refer to `./piot -h` to see proper syntax). E.g. this is how to increase log
level (`DEBUG`) for one specific command (`things`)

```
./piot --log-level DEBUG things
```

## Dump configuration

You can always check current configuration by `config` command:
```
./piot config
```

# Commands

## User profile

Command will print your current profile (email, your organizations, active
organization, etc.)

```
./piot profile
```

## Organizations

Command lists all organizations you are authorized to see:

```
./piot org

NAME          MEMBER   CURRENT   INFLUXDB
PIOT          X                  piot
JASO          X        X         jaso
```
where `MEMBER` column indicates if you are member of given org and `CURRENT`
column indicates your current organization

## Set Current Organization

Command for setting current organization, which is later used for other commands
and queries (e.g. listing of things or exports). This is how to set your current
organization to `JASO`:

```
./piot org set JASO
```

## Things

List all things assigned to your current organization:

```
./piot thing

NAME          ALIAS         TYPE/CLASS           ENABLED   LAST SEEN   VALUE
B3007-Temp                  sensor/temperature   true      44s         4.6
B3007                       device/              true      44s
B3006-Temp2   B3006-Temp2   sensor/temperature   true      42s         -13.7
B3006-Temp1   B3006-Temp1   sensor/temperature   true      42s         -15.4
B3006                       device/              true      42s
```

## Export

### Things

Export things from current org to json:

```
./piot export things
```

Export things from current org to csv:
```
./piot export things --format csv
```

### Sensors

This command exports sensor readings from Influx database. It is possible to set
time range, where default is *last 24 hours*.

Export sensors from current org to json (*last 24 hours*):
```
./piot export sensors 
```

Export sensors from current org to csv (*last 24 hours*):
```
./piot export sensors --format csv
```

Export sensors from current org to xlsx file (*last 24 hours*):
```
./piot export sensors --format xlsx -o sensors.xlsx
```

Export only sensors selected by name, for specific time interval (two days) to
csv format:
```
./piot export sensors --names B3007-Temp,B3006-Temp1 --format csv --from 2021-06-20 --to 2021-06-22
```
