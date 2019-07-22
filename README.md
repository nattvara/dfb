# Domain based Filesystem Backup

A backup approach and tool built on top of [restic](https://github.com/restic/restic) for macOS.

## About

Domain based Filesystem Backup (dfb) is backup of files organized by domains. A domain is arbitrary, but should be as narrow as possible but not to narrow. A single project such as a source code repository, documents for a project, or something longerlived such as a .photoslibrary.

Domains also include special directories such as `~/Library` and files in home directory such as `~/.zshrc` and `~/.gitconfig`.

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

dep ensure

./build.sh
```

#### Run the installer

```bash
./install.sh

dfb --version
# 1.0
```
