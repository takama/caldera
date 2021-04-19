#!/bin/bash

GO_OS=${GO_OS:-"linux"}

function detect_os {
    # Detect the OS name
    case "$(uname -s)" in
      Darwin)
        host_os=darwin
        ;;
      Linux)
        host_os=linux
        ;;
      *)
        echo "Unsupported host OS. Must be Linux or Mac OS X." >&2
        exit 1
        ;;
    esac

   GO_OS="${host_os}"
}

detect_os

old_version=$(grep 'RELEASE ?= ' ./Makefile | sed -e 's/RELEASE ?= \(.*\)/\1/g')
release_version=$(git rev-parse --abbrev-ref HEAD | grep -E "^(release|hotfix)/.+" | cut -d"/" -f2 )
if [ -z "$release_version" ]
then
  supposed_version=$(echo $old_version | awk -F. '{OFS=".";$NF = $NF + 1;} 1')
else
  supposed_version=${release_version}
fi

echo "Current version $old_version."
echo -n "Please enter new version [$supposed_version]:"
read new_version
new_version=${new_version:-$supposed_version}

if [[ $new_version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Bumped new version: $new_version"
else
    echo "Version is incorrect, please use vX.X.X format (ie: v0.17.3)"
    exit
fi

if [ "${GO_OS}" == "darwin" ]; then
    sed -i '' -e "s/\(RELEASE ?= \).*/\1$new_version/" Makefile
    sed -i '' -e "s/\(# Version\)\(\s*\).*/\1 $new_version/" docs/CHANGELOG.md
    sed -i '' -e "s/\(## Version\)\(\s*\).*/\1 $new_version/" README.md
else
    sed -i -e "s/\(RELEASE ?= \).*/\1$new_version/" Makefile
    sed -i -e "s/\(# Version\)\(\s*\).*/\1 $new_version/" docs/CHANGELOG.md
    sed -i -e "s/\(## Version\)\(\s*\).*/\1 $new_version/" README.md
fi
