# The Things Network (MQTT) Data Source for Grafana

![Grafana Data Source for The Things Network](https://lupyuen.github.io/images/grafana-flow.jpg)

Read the articles...

-   ["Grafana Data Source for The Things Network"](https://lupyuen.github.io/articles/grafana)

-   ["The Things Network on PineDio Stack BL604 RISC-V Board"](https://lupyuen.github.io/articles/ttn)

See this Twitter Thread for updates...

https://twitter.com/MisterTechBlog/status/1440459917828050946

This Grafana Data Source connects to [__The Things Network__](https://www.thethingsnetwork.org/) over MQTT to stream the received sensor data...

-   ["The Things Network: MQTT Server"](https://www.thethingsindustries.com/docs/integrations/mqtt/)

![Grafana Data Source for The Things Network](https://lupyuen.github.io/images/grafana-datasource.png)

![Grafana Dashboard for The Things Network](https://lupyuen.github.io/images/grafana-dashboard.png)

We assume that Message Payloads are encoded in [__CBOR Format__](https://en.wikipedia.org/wiki/CBOR)...

```json
{ "t": 1234 }
```

(Multiple fields are OK)

This Data Source should be located in the __Grafana Plugins Folder__...

```bash
##  For Linux: Need "sudo" to access this folder
cd /var/lib/grafana/plugins

##  For macOS: 
cd /usr/local/var/lib/grafana/plugins

##  For Windows: Need to grant "Full Control" permission to "Users" group for this folder
cd C:\Program Files\GrafanaLabs\grafana\data\plugins

##  Download source files for The Things Network Data Source
git clone --recursive https://github.com/lupyuen/the-things-network-datasource
```

To install the __Build Tools__...

-   [__Build Tools for Linux (Ubuntu)__](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894477802)

-   [__Build Tools for macOS__](https://lupyuen.github.io/articles/grafana#appendix-install-build-tools-for-macos)

-   [__Build Tools for Windows__](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894534196)

[(More details here)](https://grafana.com/tutorials/build-a-streaming-data-source-plugin/)

To __build the Data Source__...

```bash
##  Install the dependencies
cd the-things-network-datasource
yarn install

##  Build the Data Source (React + Go)
yarn build

##  If "mage" is not found, set the PATH
export PATH=$PATH:$GOPATH/bin
```

(See Build Instructions and Build Log below)

To __enable the Data Source__, edit...

```text
## For Linux:
/usr/share/grafana/conf/defaults.ini

## For macOS:
/usr/local/etc/grafana/grafana.ini

## For Windows:
C:\Program Files\GrafanaLabs\grafana\conf\defaults.ini
```

And set...

```text
[plugins]
allow_loading_unsigned_plugins = the-things-network-datasource
```

To __enable Debug Logs__, set...

```text
[log]
level = debug
```

__Configure the Data Source__ with these values from The Things Network → Application → (Your Application) → Integrations → MQTT...

```text
## Change this to your MQTT Public Address
Public Address: au1.cloud.thethings.network:1883

## Change this to your MQTT Username
Username: luppy-application@ttn

## Change this to your API Key
Password: YOUR_API_KEY

## Subscribe to all topics
Topic: all
```

To __test the MQTT Server__...

```bash
## Change au1.cloud.thethings.network to your MQTT Public Address
## Change luppy-application@ttn to your MQTT Username
mosquitto_sub -h au1.cloud.thethings.network -t "#" -u "luppy-application@ttn" -P "YOUR_API_KEY" -d
```

(See MQTT Log below)

In case of problems, check the __Grafana Log__ at...

```text
## For Linux:
/var/log/grafana/grafana.log

## For macOS:
/usr/local/var/log/grafana/grafana.log

## For Windows:
C:\Program Files\GrafanaLabs\grafana\data\log\grafana.log
```

(See Grafana Log below)

This Data Source is based on the MQTT data source for Grafana...

-   [github.com/grafana/mqtt-datasource](https://github.com/grafana/mqtt-datasource)

## Requirements

The MQTT data source has the following requirements:

- Grafana user with a server or organization administration role; refer to [Permissions](https://grafana.com/docs/grafana/latest/permissions/).

- Access to MQTT Server at The Things Network.

## Known limitations

- Only one topic is supported: "`all`"

- Message payloads must be smaller than 50 bytes [(See this)](https://github.com/lupyuen/the-things-network-datasource/blob/main/pkg/mqtt/client.go#L112-L125)

## Install the plugin

### Installation Pre-requisites

Refer to: [Building a Streaming Datasource Backend Plugin](https://grafana.com/tutorials/build-a-streaming-data-source-plugin/)

Details: [Ubuntu](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894477802) [Windows](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894534196)

### Meet compatibility requirements

__Note: Since this plugin uses the Grafana Live Streaming API, make sure to use Grafana v8.0+__
### Installation Steps

1. Clone the plugin to your Grafana plugins directory.
2. Build the plugin by running `yarn install` and then `yarn build`.

NOTE: The `yarn build` command above might fail on a non-unix-like system, like Windows, where you can try replacing the `rm -rf` command with `rimraf` in the `./package.json` file to make it work.

3. Run `mage reloadPlugin` or restart Grafana for the plugin to load.

### Verify that the plugin is installed

1. In Grafana from the left-hand menu, navigate to **Configuration** > **Data sources**.
2. From the top-right corner, click the **Add data source** button.
3. Search for "The Things Network" in the search field, and hover over "The Things Network" search result.
4. Click the **Select** button for "The Things Network".

![Grafana Data Source for The Things Network](https://lupyuen.github.io/images/grafana-datasource2.png)

## Configure the data source

[Add the Data Source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) for "The Things Network"

Configure the Data Source with the values from `The Things Network → Application → (Your Application) → Integrations → MQTT`...

#### Basic fields

| Field | Description                                        |
| ----- | -------------------------------------------------- |
| Name  | Name for this data source |
| Host  | Public Address of your MQTT Server at The Things Network |
| Port  | MQTT Port (default 1883) |

#### Authentication fields

| Field    | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| Username | Username for your MQTT Server at The Things Network |
| Password | API Key for your MQTT Server at The Things Network |

![Configuring the Grafana Data Source for The Things Network](https://lupyuen.github.io/images/grafana-config.png)

## Grafana Log

```text
C:\Program Files\GrafanaLabs\grafana\data\log\grafana.log

Subscribing to MQTT topic: #
ToFrame: topic=all
Client wants to subscribe user=1 client=14335827-eb07-4f75-890c-4df7991a031b channel=1/ds/09RFPONnk/all
Found cached channel handler channel=ds/09RFPONnk/all
Skip running new stream (already exists).features path=all
Client subscribed user=1 client=14335827-eb07-4f75-890c-4df7991a031b channel=1/ds/09RFPONnk/all

Received MQTT Message for topic v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/join
Missing or invalid payload: 
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:45:57.644920529Z",
    "join_accept": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "received_at": "2021-09-25T13:45:55.836190350Z"
    }
}

Received MQTT Message for topic v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/up
Stream MQTT Message for topic all
ToFrame: topic=all
jsonMessagesToFrame: topic=all, msg=
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:46:17.083379844Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_port": 2,
        "frm_payload": "omF0GQTSYWwZA+g=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:24:38.670302Z",
                "timestamp": 2303670303,
                "rssi": -52,
                "channel_rssi": -52,
                "snr": 11.8,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "...",
                "channel_index": 2
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "922200000",
            "timestamp": 2303670303,
            "time": "2021-09-25T14:24:38.670302Z"
        },
        "received_at": "2021-09-25T13:46:16.874919502Z",
        "consumed_airtime": "0.370688s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}
payload: [162 97 116 25 4 210 97 108 25 3 232]
CBOR decoded: map[device_id:eui-YOUR_DEVICE_EUI l:1000 t:1234]
jsonMessagesToFrame: Frame=&{Name:all Fields:[0xc0003d1b60 0xc0003d1c20 0xc0003d1bf0 0xc0003d1bc0] RefID: Meta:<nil>}
  field=&{Name:Time Labels: Config:<nil> vector:0xc0004100a0}
  field=&{Name:device_id Labels: Config:<nil> vector:0xc000410140}
  field=&{Name:l Labels: Config:<nil> vector:0xc000410120}
  field=&{Name:t Labels: Config:<nil> vector:0xc000410100}
Sending message to client for topic all
ToFrame: topic=all
jsonMessagesToFrame: topic=all, msg=
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:46:17.083379844Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_port": 2,
        "frm_payload": "omF0GQTSYWwZA+g=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:24:38.670302Z",
                "timestamp": 2303670303,
                "rssi": -52,
                "channel_rssi": -52,
                "snr": 11.8,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "...",
                "channel_index": 2
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "922200000",
            "timestamp": 2303670303,
            "time": "2021-09-25T14:24:38.670302Z"
        },
        "received_at": "2021-09-25T13:46:16.874919502Z",
        "consumed_airtime": "0.370688s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}
payload: [162 97 116 25 4 210 97 108 25 3 232]
CBOR decoded: map[device_id:eui-YOUR_DEVICE_EUI l:1000 t:1234]
jsonMessagesToFrame: Frame=&{Name:all Fields:[0xc000270ab0 0xc000270b70 0xc000270b40 0xc000270b10] RefID: Meta:<nil>}
  field=&{Name:Time Labels: Config:<nil> vector:0xc000272400}
  field=&{Name:device_id Labels: Config:<nil> vector:0xc0002724a0}
  field=&{Name:l Labels: Config:<nil> vector:0xc000272480}
  field=&{Name:t Labels: Config:<nil> vector:0xc000272460}

Received MQTT Message for topic v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/up
Stream MQTT Message for topic all
ToFrame: topic=all
jsonMessagesToFrame: topic=all, msg=
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:46:41.932243351Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_port": 2,
        "f_cnt": 1,
        "frm_payload": "omF0GQU2YWwZBEw=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:25:28.495276Z",
                "timestamp": 2328495276,
                "rssi": -51,
                "channel_rssi": -51,
                "snr": 11.5,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "..."
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "923200000",
            "timestamp": 2328495276,
            "time": "2021-09-25T14:25:28.495276Z"
        },
        "received_at": "2021-09-25T13:46:41.706545356Z",
        "consumed_airtime": "0.370688s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}
payload: [162 97 116 25 5 54 97 108 25 4 76]
CBOR decoded: map[device_id:eui-YOUR_DEVICE_EUI l:1100 t:1334]
jsonMessagesToFrame: Frame=&{Name:all Fields:[0xc0002b0ab0 0xc0002b0b70 0xc0002b0b40 0xc0002b0b10] RefID: Meta:<nil>}
  field=&{Name:Time Labels: Config:<nil> vector:0xc0002a2480}
  field=&{Name:device_id Labels: Config:<nil> vector:0xc0002a2520}
  field=&{Name:l Labels: Config:<nil> vector:0xc0002a2500}
  field=&{Name:t Labels: Config:<nil> vector:0xc0002a24e0}
Sending message to client for topic all

Received MQTT Message for topic v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/up
Missing or invalid payload: 
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:46:47.619509371Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_cnt": 2,
        "frm_payload": "Y7c=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:25:40.196422Z",
                "timestamp": 2334196423,
                "rssi": -52,
                "channel_rssi": -52,
                "snr": 11.2,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "...",
                "channel_index": 1
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "923400000",
            "timestamp": 2334196423,
            "time": "2021-09-25T14:25:40.196422Z"
        },
        "received_at": "2021-09-25T13:46:47.408241526Z",
        "consumed_airtime": "0.329728s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}

Received MQTT Message for topic v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/up
Stream MQTT Message for topic all
ToFrame: topic=all
jsonMessagesToFrame: topic=all, msg=
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:47:12.461507414Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_port": 2,
        "f_cnt": 3,
        "frm_payload": "omF0GQTcYWwZA/I=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:26:29.039893Z",
                "timestamp": 2359039893,
                "rssi": -52,
                "channel_rssi": -52,
                "snr": 11.5,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "...",
                "channel_index": 4
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "922600000",
            "timestamp": 2359039893,
            "time": "2021-09-25T14:26:29.039893Z"
        },
        "received_at": "2021-09-25T13:47:12.247172827Z",
        "consumed_airtime": "0.370688s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}
payload: [162 97 116 25 4 220 97 108 25 3 242]
CBOR decoded: map[device_id:eui-YOUR_DEVICE_EUI l:1010 t:1244]
jsonMessagesToFrame: Frame=&{Name:all Fields:[0xc0001dff20 0xc000286030 0xc0001dffb0 0xc000286540] RefID: Meta:<nil>}
  field=&{Name:Time Labels: Config:<nil> vector:0xc000080800}
  field=&{Name:device_id Labels: Config:<nil> vector:0xc000080880}
  field=&{Name:l Labels: Config:<nil> vector:0xc000080860}
  field=&{Name:t Labels: Config:<nil> vector:0xc0000808a0}
Sending message to client for topic all
```

## MQTT Log

Monitor MQTT Server at The Things Network with `mosquitto_sub`...

```text
"c:\Program Files\Mosquitto\mosquitto_sub" -h au1.cloud.thethings.network -t "#" -u "luppy-application@ttn" -P "YOUR_API_KEY" -d

Client (null) sending CONNECT
Client (null) received CONNACK (0)
Client (null) sending SUBSCRIBE (Mid: 1, Topic: #, QoS: 0, Options: 0x00)
Client (null) received SUBACK
Subscribed (mid: 1): 0
Client (null) sending PINGREQ
Client (null) received PINGRESP
```

Join Request...

```text
Client (null) received PUBLISH (d0, q0, r0, m0, 'v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/join', ... (691 bytes))
```

```json
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:45:57.644920529Z",
    "join_accept": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "received_at": "2021-09-25T13:45:55.836190350Z"
    }
}
```

Send Data...

```text
Client (null) received PUBLISH (d0, q0, r0, m0, 'v3/luppy-application@ttn/devices/eui-YOUR_DEVICE_EUI/up', ... (1453 bytes))
```

```json
{
    "end_device_ids": {
        "device_id": "eui-YOUR_DEVICE_EUI",
        "application_ids": {
            "application_id": "luppy-application"
        },
        "dev_eui": "YOUR_DEVICE_EUI",
        "join_eui": "0000000000000000",
        "dev_addr": "YOUR_DEV_ADDR"
    },
    "correlation_ids": [
        "as:up:...",
        "gs:conn:...",
        "gs:up:host:...",
        "gs:uplink:...",
        "ns:uplink:...",
        "rpc:/ttn.lorawan.v3.GsNs/HandleUplink:...",
        "rpc:/ttn.lorawan.v3.NsAs/HandleUplink:..."
    ],
    "received_at": "2021-09-25T13:46:17.083379844Z",
    "uplink_message": {
        "session_key_id": "YOUR_SESSION_KEY_ID",
        "f_port": 2,
        "frm_payload": "omF0GQTSYWwZA+g=",
        "rx_metadata": [
            {
                "gateway_ids": {
                    "gateway_id": "luppy-wisgate-rak7248",
                    "eui": "YOUR_EUI"
                },
                "time": "2021-09-25T14:24:38.670302Z",
                "timestamp": 2303670303,
                "rssi": -52,
                "channel_rssi": -52,
                "snr": 11.8,
                "location": {
                    "latitude": 1.27125,
                    "longitude": 103.80795,
                    "altitude": 70,
                    "source": "SOURCE_REGISTRY"
                },
                "uplink_token": "...",
                "channel_index": 2
            }
        ],
        "settings": {
            "data_rate": {
                "lora": {
                    "bandwidth": 125000,
                    "spreading_factor": 10
                }
            },
            "data_rate_index": 2,
            "coding_rate": "4/5",
            "frequency": "922200000",
            "timestamp": 2303670303,
            "time": "2021-09-25T14:24:38.670302Z"
        },
        "received_at": "2021-09-25T13:46:16.874919502Z",
        "consumed_airtime": "0.370688s",
        "network_ids": {
            "net_id": "000013",
            "tenant_id": "ttn",
            "cluster_id": "ttn-au1"
        }
    }
}
```

## Build Log

```text
C:\Program Files\GrafanaLabs\grafana\data\plugins\the-things-network-datasource>yarn build 
yarn run v1.22.11
$ rm -rf dist && grafana-toolkit plugin:build && mage build:backend
  Using Node.js v14.17.6
  Using @grafana/toolkit v8.0.0-beta.3
√ Preparing
√ Linting
ts-jest[config] (WARN) The option `tsConfig` is deprecated and will be removed in ts-jest 27, use `tsconfig` instead
 PASS  src/handleEvent.test.ts

Test Suites: 1 passed, 1 total
Tests:       2 passed, 2 total
Snapshots:   2 passed, 2 total
Time:        2.521 s
Ran all test suites with tests matching "".
√ Running tests
\ Compiling...  Starting type checking service...
  Using 1 worker with 2048MB memory limit
/ Compiling...  
   Hash: 6b5a018b08c2ac55e195
  Version: webpack 4.41.5
  Time: 9050ms
  Built at: 09/21/2021 3:35:48 PM
                  Asset       Size  Chunks                   Chunk Names
           CHANGELOG.md   58 bytes          [emitted]        
                LICENSE   11.3 KiB          [emitted]        
              README.md   5.22 KiB          [emitted]        
           img/mqtt.svg   1.33 KiB          [emitted]        
              module.js    3.8 KiB       0  [emitted]        module
  module.js.LICENSE.txt  808 bytes          [emitted]        
          module.js.map   24.4 KiB       0  [emitted] [dev]  module
            plugin.json   1.11 KiB          [emitted]        
  Entrypoint module = module.js module.js.map
  [0] external "react" 42 bytes {0} [built]
  [1] external "@grafana/ui" 42 bytes {0} [built]
  [2] external "lodash" 42 bytes {0} [built]
  [3] external "@grafana/data" 42 bytes {0} [built]
  [4] external "@grafana/runtime" 42 bytes {0} [built]
  [5] ./module.ts + 5 modules 14.5 KiB {0} [built]
      | ./module.ts 296 bytes [built]
      | ./datasource.ts 352 bytes [built]
      | ./ConfigEditor.tsx 2.78 KiB [built]
      | ./QueryEditor.tsx 658 bytes [built]
      | ../node_modules/tslib/tslib.es6.js 10 KiB [built]
      | ./handleEvent.ts 395 bytes [built] 
  
√ Compiling...
Done in 37.51s.
```
