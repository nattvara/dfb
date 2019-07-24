# Domain based Filesystem Backup

A backup approach and tool built on top of [restic](https://github.com/restic/restic) for macOS.

## About

Domain based Filesystem Backup (dfb) is backup of files organized by domains. A domain is arbitrary, but should be as narrow as possible but not to narrow. A single project such as a source code repository, documents for a project, or something longerlived such as a `.photoslibrary`.

Domains also include special directories such as `~/Library` and files in home directory such as `~/.zshrc` and `~/.gitconfig`.

The purpose of dfb is to allow for easier management of backups, better tools for monitoring backups that run in the background, easier recovery of backed up files, and easier retrival of useful stats about repos than a standalone restic installation provides.

## Installation

### Build and install locally

#### Install [restic](https://github.com/restic/restic)

```bash
brew install restic
```

#### Install FUSE for macOS

Either from [official website (recommended)](https://github.com/osxfuse/osxfuse) or using homebrew

```bash
brew update
brew tap homebrew/cask
brew cask install osxfuse
```

#### Install GNU Grep

```bash
brew install grep
```

#### Install go and dep

```bash
brew install go
brew install dep
```

#### Download and build dfb and dependencies

```bash
git clone https://github.com/nattvara/dfb.git
cd dfb

if [ ! -L "$GOPATH/src/dfb" ]; then ln -s "$(pwd)" "$GOPATH/src/dfb"; else echo "already exists"; fi

cd "$GOPATH/src/dfb"

go get fyne.io/fyne
dep ensure

./build.sh
```

#### Run the installer

```bash
./install.sh

dfb --version
# 1.0
```

## Usage

### Structuring the filesystem

The purpose of dfb is to keep domains separate and allow for backup, recovery and deletion of backups easy and isolated. Hence the filesystem must be organized in such a way.

```console
~
├── some-project
├── someorg-some-project -> /Volumes/some_disk/someorg-some-project
├── someorg-some-other-project
├── Desktop     # directory that shouldn't be backed up
├── Downloads   # directory that shouldn't be backed up
├── Library     # directory that should be backed up, but not renamed
└── ...
```

### Groups

Domains should be organized by groups. A group contians a number of domains, and restic repositories to backup those domains to.

### Repositories

Repositories are created with restic.

```bash
restic -r [REPO] init
```

[See the restic docs.](https://restic.readthedocs.io/)

### Availible commands

```console
$ dfb

_________/\\\__________/\\\\\___/\\\________
 ________\/\\\________/\\\///___\/\\\________
  ________\/\\\_______/\\\_______\/\\\________
   ________\/\\\____/\\\\\\\\\____\/\\\________
    ___/\\\\\\\\\___\////\\\//_____\/\\\\\\\\\__
     __/\\\////\\\______\/\\\_______\/\\\////\\\_
      _\/\\\__\/\\\______\/\\\_______\/\\\__\/\\\_
       _\//\\\\\\\\\______\/\\\_______\/\\\\\\\\\\_
        __\/////////_______\///________\//////////__

Domain based Filesystem Backup.

Usage:
  dfb <command> <subcommand> [parameters]

Available Commands:
  groups      Group commands.
  domain      Domain commands.
  backup      Backup a group of domains to a repo.
  recover     Mount backed up versions of domains for recovery.
  fsd         Control the filesystem agent.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
```

## License

MIT © Ludwig Kristoffersson

See [LICENSE file](LICENSE) for more information.
