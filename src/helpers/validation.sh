#
# Summary: Validation functions
#

validate_group() {
    if [ "$1" == "" ]; then
        echo "please provide a group."
        exit 1
    fi

    if [ ! -d "$DFB_PATH/$1" ]; then
        echo "please provide a valid group."
        exit 1
    fi
}

validate_domain() {
    if [ "$1" == "" ]; then
        echo "please provide a group."
        exit 1
    fi

    if [ ! -d "$DFB_PATH/$1" ]; then
        echo "please provide a valid group."
        exit 1
    fi

    if [ "$2" == "" ]; then
        echo "please provide a domain."
        exit 1
    fi

    if [ ! -f "$DFB_PATH/$group/domains/$domain" ]; then
        echo "please provide a valid domain."
        exit 1
    fi
}

validate_repo() {
    if [ ! -f "$DFB_PATH/$1/repos/$2" ]; then
        echo "please provide a valid repo"
        exit 1
    fi
}
