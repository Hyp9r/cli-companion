# cli-companion

Command line tool that helps with boring tasks like creating database models and sql queries for golang

## Download

Download .deb package from [cli-companion.deb](https://github.com/Hyp9r/cli-companion/releases/tag/v0.0.1)

## Installation

```
sudo dpkg -i cli-companion_<replace with version>-<iteration for version>_amd64.deb
```
Example with version and iteration number
```
sudo dpkg -i cli-companion_0.0.1-1_amd64.deb
```

Add companion to PATH in your bashrc or zshrc file
```
export PATH=$PATH:/usr/bin/cli-companion
```

## Useage

In your project root directory in terminal just write ```companion``` to start the cli

Generating models and repositories with ```companion -operation make -namespace model```

## License

[GPL](https://choosealicense.com/licenses/gpl-3.0/)