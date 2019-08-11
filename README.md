# Domain based Filesystem Backup

A backup approach and tool for macOS built on top of [restic](https://github.com/restic/restic).

## About

Domain based Filesystem Backup (dfb) is backup of files, organized by domains. A domain is arbitrary, but should be as narrow as possible, while not to narrow. A single project such as a source code repository, documents for a project, or something longerlived such as a `.photoslibrary`.

Domains can also include special directories such as `~/Library` and files in home directory such as `~/.zshrc` and `~/.gitconfig`.

The purpose of dfb is to provide handy monitoring tools for backups that run in the background, convenient recovery of backed up files, management of backup locations and exceptions for individual domains, and easier retrival of useful stats about the backups than a standalone restic installation provides.

The tool is written in bash with the parts that require higher performance written in go.

## Installation

### Build and install locally

#### Install dependencies

```bash
brew install restic go expect jq grep coreutils terminal-notifier
```

#### Install FUSE for macOS

Either from [official website (recommended)](https://github.com/osxfuse/osxfuse) or using homebrew

```bash
brew update
brew tap homebrew/cask
brew cask install osxfuse
```

#### Install fyne

```bash
go get fyne.io/fyne/cmd/fyne
```

#### Download and build dfb

```bash
git clone https://github.com/nattvara/dfb.git
cd dfb

./build.sh
```

#### Run the installer

```bash
./install.sh

dfb --version
# 1.1
```

## Usage

### Structuring the filesystem

One of the purposes of dfb is to keep snapshots of domains separate, so that recovery and deletion of backups are easy and isolated. Hence the filesystem must be organized in such a way.

```console
~
├── some-project
├── someorg-some-project -> /Volumes/some_disk/someorg-some-project
├── someorg-some-other-project
├── Library
├── Desktop     # directory that shouldn't be backed up
├── Downloads   # directory that shouldn't be backed up
└── ...
```

### Repositories

Repositories are created with restic.

```bash
restic -r [RESTIC REPO] init
```

[See the restic docs.](https://restic.readthedocs.io/)

### Groups

Domains are backed up in groups. A group contians a number of domains, and restic repositories to backup those domains to.

Creating groups is done with the `groups` command.

```console
$ dfb groups add
# Enter name of new group: demo
```

Groups can have one or more repos to back up to.

```console
$ dfb groups add-repo demo
# Enter name of repo: demo-repo
# Enter repo path: [RESTIC REPO]
```

#### Available subcommands for the `groups` command

```console
$ dfb groups --help
Domain groups.

A group contians a number of domains, and restic repositories
to backup those domains to.

Usage:
  dfb groups <subcommand> [parameters]

Available Commands:
  ls        List groups.
  add       Add new group.
  repos     List restic repos for a group.
  add-repo  Add restic repo for a group.

Options:
  -h --help     Show this screen.
```

### Domains

Adding domains to a group for backup is done with the `domains` command.

```bash
dfb domains add demo ~/demo-some-project
dfb domains add demo ~/demo-some-other-project
```

#### Symlinked domains

Domains can have their real source at another location than the `$HOME` directory. An example of this would be storing a domain on an external drive.

```bash
dfb domains add demo ~/demo-a-symlinked-domain /Volumes/[SOME VOLUME]/demo-a-symlinked-domain
```

In the background dfb:s filesystem agent will detect when the volume is availible and create a symkink to the real directory in the users `$HOME` directory.

#### Available subcommands for the `domains` command

```console
$ dfb domains --help
Domains to backup.

A domain is a directory in the home directory to backup,
this could be a symlink to some other directory on another
volume.

Usage:
  dfb domains <subcommand> [parameters]

Available Commands:
  ls        List domains.
  add       Add new domain.
  rm        Remove a domain.

Options:
  -h --help     Show this screen.
```

### Backup

Backups are done with the `backup` command.

Backup will launch a macOS dialogue to enter the password and subsequently take a snapshot of all domains that are *availible* (see [symlinked domains](#symlinked-domains)).

```console
$ dfb backup demo demo-repo
  backing up demo-a-symlinked-domain              100% ⏱  2.3s 💾 309.1 MiB 📊 gathering stats... done.
  backing up demo-some-other-project              100% ⏱  0.7s 💾 99.3 MiB 📊 gathering stats... done.
  backing up demo-some-project                    100% ⏱  0.3s 💾 0 B 📊 gathering stats... done.


  gathering stats for repo demo-repo... done
```

#### `--gui`

With the `---gui` flag, progress will be displayed in a gui window. This is useful if the backup is started from a cron job or similar.

![Backup with --gui flag](docs/images/progress-gui.png)

#### `--confirm`

With the `--confirm` flag a dialogue will promt the user for action before backup starts. Also helpful if the backup is started by a cron job.

![Backup with --confirm flag](docs/images/confirm-dialogue.png)

#### Available subcommands for the `backup` command

```console
$ dfb backup --help
Backup a group of domains.

Usage:
  dfb backup [group] [repo]

Options:
  --gui         Show progress in a graphical user interface.
  --confirm     Show a dialogue that the user have to confirm for backup to start. Useful if backup is started by a cron job.
  --force       Force backup to start, even if dfb is locked.
  -h --help     Show this screen.
```

### Recovering files

The `recover` command will run the `restic mount` command that mounts a read-only filesystem with provided repo. In the background dfb will create aliases to a domains backups under the path `path/to/domain/__recover__` like the following.

```console
demo-some-project
├── __recover__ -> ~/.dfb/demo/mountpoint/tags/demo-some-project
│   ├── 2019-01-04T10:00:41+02:00
│   │   └── some-file.txt
│   ├── 2019-03-04T12:43:31+02:00
│   │   └── some-file.txt
│   ├── 2019-04-04T10:30:02+02:00
│   │   └── some-file.txt
│   ├── 2019-05-04T02:12:39+02:00
│   │   └── some-file.txt
│   ├── 2019-07-04T07:12:30+02:00
│   │   ├── some-file.txt
│   │   └── some-other-file.txt
│   ├── 2019-07-04T01:10:31+02:00
│   │   ├── some-file.txt
│   │   └── some-other-file.txt
│   ├── 2019-08-04T17:21:59+02:00
│   │   ├── some-file.txt
│   │   └── some-other-file.txt
│   └── latest -> 2019-08-04T17:21:59+02:00  [recursive, not followed]
├── some-file.txt
└── some-other-file.txt
14 directories, 13 files
```

#### Available subcommands for the `recover` command

```console
$ dfb recover --help
Recover files from a group.

The recover command will mount given restic repo for a group, once
mounted the dfb agent will create a __recover__ directory under each
domain that contains earlier versions of all the backed up files in the repo.

Usage:
  dfb recover [group] [repo]

Options:
  --force       Force recover to start, even if dfb is locked.
  -h --help     Show this screen.
```

### Stats

The `stats` command offers a few helpful metrics to gain insight into how much space your backups occupy, which domains take up the most space or if time taken for backup is increasing/decreasing.

The following can be useful if restic repositories are stored in an environment where space is not free, on AWS S3 for instance.

```bash
dfb stats demo demo-repo repo-disk-space --time-unit days --time-length 34
```

![Example usage of the stats command](docs/images/stats-repo-disk-space-example.png)

#### Full list of options for the `stats` command

```console
Usage:
  stats [group] [repo] [metric] [flags]

Flags:
  -a, --aggregator string   aggregation method to use for a metric
  -d, --domain string       which domain to use for metric, not availiable for all metrics, optional/required for some metrics
  -h, --help                help for stats
      --list-aggregators    list availiable aggregators
      --list-metrics        list availiable metrics
      --list-time-units     list availiable time units
  -o, --output string       output path for png image of metric (default "/tmp/dfb-metric.png")
  -l, --time-length int     how many time-units of history should be included (default 7)
  -u, --time-unit string    time unit to use for metric (default "days")
```

### All availible commands

```console
$ dfb --help

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
  stats       Make a chart for a backup metric.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
```

## License

__MIT License__, see [LICENSE file](LICENSE).
