#
# Summary: Main script for dfb
#
# Domain based Filesystem Backup (dfb) is backup of files
# organized by domains. A domain is arbitrary, but should be
# as narrow as possible but not to narrow. A single project
# such as a source code repository, documents for a project,
# or something longerlived such as a .photoslibrary. Domains
# also include special directories such as ~/Library and
# files in home directory such as ~/.zshrc and ~/.gitconfig.

DFB_PATH="$HOME/.dfb"
PROGRAM=$(basename "${0}")
VERSION=1.0

main() {
    if [[ "${1:-}" =~ ^-h|--help$  ]]
    then
        print_main_help
    elif [[ "${1:-}" =~ ^-v|--version$  ]]
    then
        print_version
    elif [ "${1:-}" == "groups" ]
    then
        groups "$@"
    elif [ "${1:-}" == "domains" ]
    then
        domains "$@"
    elif [ "${1:-}" == "backup" ]
    then
        backup "$@"
    else
        print_main_help
    fi
}

print_main_help() {
    cat <<HEREDOC

_________/\\\\\\__________/\\\\\\\\\\___/\\\\\\________
 ________\\/\\\\\\________/\\\\\\///___\\/\\\\\\________
  ________\\/\\\\\\_______/\\\\\\_______\\/\\\\\\________
   ________\\/\\\\\\____/\\\\\\\\\\\\\\\\\\____\\/\\\\\\________
    ___/\\\\\\\\\\\\\\\\\\___\\////\\\\\\//_____\\/\\\\\\\\\\\\\\\\\\__
     __/\\\\\\////\\\\\\______\\/\\\\\\_______\\/\\\\\\////\\\\\\_
      _\\/\\\\\\__\\/\\\\\\______\\/\\\\\\_______\\/\\\\\\__\\/\\\\\\_
       _\\//\\\\\\\\\\\\\\\\\\______\\/\\\\\\_______\\/\\\\\\\\\\\\\\\\\\\\_
        __\\/////////_______\\///________\\//////////__

Domain based Filesystem Backup.

Usage:
  ${PROGRAM} <command> <subcommand> [parameters]

Available Commands:
  groups      Group commands.
  domain      Domain commands.
  backup      Backup a group of domains to a repo.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
HEREDOC
}

verify_env() {
    if [[ ! "$OSTYPE" == "darwin"* ]]; then
        echo "dfb is only availible for macOS"
        exit 1
    fi

    if [ "$(command -v restic)" == "" ]; then
        echo "restic is not installed, visit https://github.com/restic/restic"
        exit 1
    fi

    if [ ! -d "/Library/Filesystems/osxfuse.fs" ]; then
        echo "FUSE for macOS is not installed, visit https://github.com/osxfuse/osxfuse"
        exit 1
    fi

    if [ "$(command -v ggrep)" == "" ]; then
        echo "GNU grep is not installed, run: brew install grep"
        exit 1
    fi

    if [ "$(command -v dfb-progress-parser)" == "" ]; then
        echo "dfb-progress-parser, view installation instructions in README.md that was distributed with this software"
        exit 1
    fi

    if [ ! -d "$DFB_PATH" ]; then
        echo "creating dfb root directory at $DFB_PATH"
        mkdir "$DFB_PATH"
    fi
}

print_version() {
    echo $VERSION
}
