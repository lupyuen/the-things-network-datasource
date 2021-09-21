# The Things Network (MQTT) Data Source for Grafana

Read the article...

-   ["The Things Network on PineDio Stack BL604 RISC-V Board"](https://lupyuen.github.io/articles/ttn)

This Grafana Data Source connects to [__The Things Network__](https://www.thethingsnetwork.org/) over MQTT to stream the received sensor data...

-   ["The Things Network: MQTT Server"](https://www.thethingsindustries.com/docs/integrations/mqtt/)

We assume that Message Payloads are encoded in [__CBOR Format__](https://en.wikipedia.org/wiki/CBOR)...

```json
{ "t": 1745 }
```

This Data Source should be located in the __Grafana Plugins Folder__...

```text
C:\Program Files\GrafanaLabs\grafana\data\plugins\the-things-network-datasource
```

To __build the Data Source__...

```bash
yarn install
yarn build
```

(More instructions below)

To __enable the Data Source__, edit...

```text
C:\Program Files\GrafanaLabs\grafana\conf\defaults.ini
```

And set...

```text
[plugins]
allow_loading_unsigned_plugins = the-things-network-datasource
```

__Configure the Data Source__ with these values from The Things Network -> Application -> (Your Application) -> Integrations -> MQTT...

```text
## Change this to your MQTT Public Address
Public Address: au1.cloud.thethings.network:1883

## Change this to your MQTT Username
Username: luppy-application@ttn

## Change this to your API Key
Password: <YOUR_API_KEY>

## For all topics
Topic: #

## For a specific device
Topic: v3/{application id}@{tenant id}/devices/{device id}/up
```

Based on the MQTT data source for Grafana...

-   [github.com/grafana/mqtt-datasource](https://github.com/grafana/mqtt-datasource)

## TODO: Requirements

The MQTT data source has the following requirements:

- Grafana user with a server or organization administration role; refer to [Permissions](https://grafana.com/docs/grafana/latest/permissions/).
- Access to a MQTT broker.

## TODO: Known limitations

- The plugin currently does not support all of the MQTT CONNECT packet options.
- The plugin currently does not support TLS.
- Including multiple topics in a panel is not yet well supported.
- This plugin automatically supports topics publishing very simple JSON formatted messages. Note that only the following structure is supported as of now:
```
{
    'value1': 1.0,
    'value2': 2,
    'value3': 3.33,
    ...
}
```
We do plan to support more complex JSON data structures in the upcoming releases. Contributions are highly encouraged!
- This plugin currently attaches timestamps to the messages when they are received, so there is no way to have custom timestamp for messages.

## TODO: Install the plugin

### TODO: Installation Pre-requisites

Refer to: [Building a Streaming Datasource Backend Plugin](https://grafana.com/tutorials/build-a-streaming-data-source-plugin/)

Details: [Ubuntu](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894477802) [Windows](https://github.com/grafana/mqtt-datasource/issues/15#issuecomment-894534196)

### TODO: Meet compatibility requirements

This plugin currently supports MQTT v3.1.x.

__Note: Since this plugin uses the Grafana Live Streaming API, make sure to use Grafana v8.0+__
### Installation Steps

1. Clone the plugin to your Grafana plugins directory.
2. Build the plugin by running `yarn install` and then `yarn build`.

NOTE: The `yarn build` command above might fail on a non-unix-like system, like Windows, where you can try replacing the `rm -rf` command with `rimraf` in the `./package.json` file to make it work.

3. Run `mage reloadPlugin` or restart Grafana for the plugin to load.

### TODO: Verify that the plugin is installed

1. In Grafana from the left-hand menu, navigate to **Configuration** > **Data sources**.
2. From the top-right corner, click the **Add data source** button.
3. Search for `MQTT` in the search field, and hover over the MQTT search result.
4. Click the **Select** button for MQTT.

## TODO: Configure the data source

[Add a data source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/) by filling in the following fields:

#### TODO: Basic fields

| Field | Description                                        |
| ----- | -------------------------------------------------- |
| Name  | A name for this particular AppDynamics data source |
| Host  | The hostname or IP of the MQTT Broker              |
| Port  | The port used by the MQTT Broker (default 1883)    |

#### TODO: Authentication fields

| Field    | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| Username | (Optional) The username to use when connecting to the MQTT broker |
| Password | (Optional) The password to use when connecting to the MQTT broker |

## TODO: Query the data source

The query editor allows you to specify which MQTT topics the panel will subscribe to. Refer to the [MQTT v3.1.1 specification](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718106)
for more information about valid topic names and filters.

![mqtt dashboard](./test_broker.gif)
