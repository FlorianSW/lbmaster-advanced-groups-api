# LBmaster Advanced Groups API

This is a little tool that can be run next to the Advanced Groups Mod from
LBmaster (https://lbmaster.de/product.php?id=4). It simply offers a basic API for the configuration to manage dynamic
parts from somewhere else.

Currently, only prefix Groups management is implemented in the API, which is the most dynamic part of the config.

## Prerequisites

This tool does work only, if:

- You have a dedicated server box, where you can run executables on, Game-Service-Providers (GSP, like GTX, Nitrado or
  so) will not work
- You have installed and configured the Advanced Groups Mod already (that's out of the scope of this Readme)
- You have some basic understanding of Server maintenance

## Setup

Basically, it's as easy as that:

- Download the latest version of this tool from
  the [automated builds](https://github.com/FlorianSW/lbmaster-advanced-groups-api/actions/workflows/build.yml) (take
  the most-top green build) or from the [releases page](https://github.com/FlorianSW/lbmaster-advanced-groups-api/releases)
    - There is a windows and linux binary, take the one for your platform
- Copy the binary to some location on your server (the server where your DayZ server runs as well)
- Start the executable once (if a window opens, close that for now)
- A new file, `config.json`, should have been created next to the executable in the same directory; open it in your
  favorite text editor
- Change the following config, according to your setup:
    - the `port` -> it defaults to `8080`, which is fine. If the port, however, is already used, choose some other
      port (e.g. when you run multiple instances of this tool)
    - the `advanced_groups_config_path` -> the absolute path to the Advanced Groups config file (`\` needs to be escaped
      with another backslash, e.g. a path like `D:\some\path\Config.json` becomes `D:\\some\\path\\Config.json`)
- Start the tool again by double-clicking it

That's it, the tool should start and be ready to serve your requests. You might need to change your firewall or router
settings to either forward or whitelist traffic for the port you just configured in the `cconfig.json` file.

## Usage

The tool offers a REST-like API, which allows you to modify some parts of the config of Advanced Groups.

### Authentication

The tool uses a pre-configured API key to authenticate and authorize requests, which can be found in the `config.json`.
Everyone who knows this API key can issue requests to the tool, and they will be fulfilled.

In order to make an authenticated request, add the `Authorization` header to your HTTP requests with the `Bearer`
scheme, e.g.

```
Authorization: Bearer abc123
```

Given the API key is `abc123`.

### Available endpoints

The following list will give you an overview of available endpoints offered by this tool:

#### `GET /api/prefixGroups`

Lists available prefix groups.

Example response:

```json
[
  {
    "index": 0,
    "prefix": "[VIP] "
  }
]
```

#### `GET /api/prefixGroups/<index>`

Lists members of that prefix group.

_Parameters_:

- <index>: The index of the prefix group to query, e.g. `0`

Example request:
`GET /api/prefixGroups/0`

Example response:

```json
[
  "76561111111111111",
  "76561122222222222"
]
```

#### `PUT /api/prefixGroups/<index>/<steamUID>`

Adds the requested Steam UID to the given prefix Group.

_Parameters_:

- <index>: The index of the prefix group to query, e.g. `0`
- <steamUID>: The steam UID of the profile you want to add to the prefix group

Example request:
`GET /api/prefixGroups/0/76561122222222222`

Example response:

The response is always empty. A successful request is indicated by a `204` status code.

#### `DELETE /api/prefixGroups/<index>/<steamUID>`

Removes the requested Steam UID from the given prefix Group.

_Parameters_:

- <index>: The index of the prefix group to query, e.g. `0`
- <steamUID>: The steam UID of the profile you want to remove from the prefix group

Example request:
`DELETE /api/prefixGroups/0/76561122222222222`

Example response:

The response is always empty. A successful request is indicated by a `204` status code.

#### `GET /api/prefixGroups/<index>/<steamUID>`

Can be used to check if a specific Steam UID is member of the given prefix Group.

_Parameters_:

- <index>: The index of the prefix group to query, e.g. `0`
- <steamUID>: The steam UID of the profile you want to check the membership of

Example request:
`GET /api/prefixGroups/0/76561122222222222`

Example response:

The response is always empty.
If the steam UID _is_ member of the prefix group, the status code will be 204.
If the steam UID _is not_ member of the prefix group, the status code will be 404.
