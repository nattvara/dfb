# Domain based Filesystem Backup

A backup approach and tool built on top of [restic](https://github.com/restic/restic)

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

#### Download and build dfb dependencies

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
