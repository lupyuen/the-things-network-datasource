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

__Configure the Data Source__ with these values from The Things Network → Application → (Your Application) → Integrations → MQTT...

```text
## Change this to your MQTT Public Address
Public Address: au1.cloud.thethings.network:1883

## Change this to your MQTT Username
Username: luppy-application@ttn

## Change this to your API Key
Password: <YOUR_API_KEY>

## Subscribe to all topics
Topic: #

## Subscribe to messages for a specific device
## Change luppy-application@ttn to your MQTT Username
Topic: v3/luppy-application@ttn/devices/{device id}/up
```

Based on the MQTT data source for Grafana...

-   [github.com/grafana/mqtt-datasource](https://github.com/grafana/mqtt-datasource)

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
