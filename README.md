# watermint toolbox

[![CircleCI](https://circleci.com/gh/watermint/toolbox.svg?style=svg)](https://circleci.com/gh/watermint/toolbox)
[![Coverage Status](https://coveralls.io/repos/github/watermint/toolbox/badge.svg)](https://coveralls.io/github/watermint/toolbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/watermint/toolbox)](https://goreportcard.com/report/github.com/watermint/toolbox)

Tools for Dropbox and Dropbox Business.

# Licensing & Disclaimers

watermint toolbox is licensed under the MIT license. Please see LICENSE.md or LICENSE.txt for more detail.

Please carefully note:

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.

# Usage

`tbx` have various features. Run without an option for a list of supported commands and options.
You can see available commands and options by running executable without arguments like below.

```bash
% ./tbx
toolbox 52.2.43
© 2016-2019 Takayuki Okazaki
Licensed under open source licenses. Use the `license` command for more detail.


Tools for Dropbox and Dropbox Business

Usage:
./tbx  [command]

Available commands:
   group         Group management (Dropbox Business)
   license       Show license information
   member        Team member management (Dropbox Business)
   sharedfolder  Shared folder
   sharedlink    Shared Link of Personal account
   team          Dropbox Business Team
   teamfolder    Team folder management (Dropbox Business)
   web           Launch web console
```

## Commands

| command                      | description                                                  |
|------------------------------|--------------------------------------------------------------|
| `file compare`               | Compare files between two accounts                           |
| `file copy`                  | Copy files                                                   |
| `file list`                  | List files/folders                                           |
| `file metadata`              | Report metadata for a file or folder                         |
| `file mirror`                | Mirror files/folders into another account                    |
| `file move`                  | Copy files                                                   |
| `file save`                  | Save the data from a specified URL into a file               |
| `group list`                 | List group(s)                                                |
| `group member add`           | Add members into existing groups                             |
| `group member list`          | List members of groups                                       |
| `group remove`               | Remove group                                                 |
| `license`                    | Show license information                                     |
| `member detach`              | Convert Dropbox Business accounts to Basic account           |
| `member invite`              | Invite member(s)                                             |
| `member list`                | List team member(s)                                          |
| `member mirror files`        | Mirror member files                                          |
| `member quota update`        | Update member storage quota                                  |
| `member remove`              | Remove the member from the team                              |
| `member sync`                | Sync member information with provided csv                    |
| `member update email`        | Update member email address                                  |
| `sharedfolder list`          | List shared folder(s)                                        |
| `sharedfolder member list`   | List shared folder member(s)                                 |
| `sharedlink create`          | Create shared link                                           |
| `sharedlink list`            | List of shared link(s)                                       |
| `sharedlink remove`          | Remove shared link                                           |
| `team audit events`          | Export activity logs                                         |
| `team audit sharing`         | Export all sharing information across team                   |
| `team device list`           | List devices or web sessions of the team                     |
| `team device unlink`         | Unlink device                                                |
| `team feature`               | Team feature                                                 |
| `team info`                  | Team information                                             |
| `team linkedapp list`        | List linked applications                                     |
| `team namespace file list`   | List files/folders in all namespaces of the team             |
| `team namespace file size`   | Calculate size of namespaces                                 |
| `team namespace list`        | List all namespaces of the team                              |
| `team namespace member list` | List all namespace members of the team                       |
| `team sharedlink cap expiry` | Force expiration date of public shared links within the team |
| `team sharedlink list`       | List of shared link(s)                                       |
| `teamfolder archive`         | Archive team folder(s)                                       |
| `teamfolder file list`       | List files/folders in all team folders of the team           |
| `teamfolder list`            | List team folder(s)                                          |
| `teamfolder mirror`          | Mirror team folders into another team                        |
| `teamfolder permdelete`      | Permanently delete team folder(s)                            |
| `teamfolder size`            | Calculate size of team folder                                |
| `web`                        | Launch web console (experimental)                            |

## Authentication

If an executable contains registered application keys, then the executable will ask you an authentication to your Dropbox account or a team.
Please open the provided URL, then paste authorisation code.

```
toolbox 52.2.43
© 2016-2019 Takayuki Okazaki
Licensed under open source licenses. Use the `license` command for more detail.

Testing network connection...
Done

1. Visit the URL for the auth dialog:

https://www.dropbox.com/oauth2/authorize?client_id=xxxxxxxxxxxxxxx&response_type=code&state=xxxxxxxx

2. Click 'Allow' (you might have to login first):
3. Copy the authorisation code:
Enter the authorisation code
```

The executable store tokens at the file under folder `$HOME/.toolbox/secrets/(HASH).secret`. If you don't want to store tokens into the file, then please specify option `-secure`.

## Proxy

The executable automatically detects your proxy configuration from the environment. However, if you got an error or you want to specify explicitly, please add `-proxy` option, like `-proxy hostname:port`.
Currently, the executable doesn't support proxies which require authentication.
