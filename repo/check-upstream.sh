#!/usr/bin/env sh

for commit in $(git rev-list origin/master..HEAD); do
    git show --pretty="format:" --name-only ${commit} > /tmp/files
    if [ -n "$(grep /packages/upstream/ /tmp/files)" ]; then
        if [ \
            "$(git log --format=%s ${commit}^..${commit})" != \
            "AUTO: Update upstream packages" ];
        then
            echo "ERROR: ${commit} manually updates"
            echo "repo/[darwin|win32]/packages/packages/upstream. Instead, please run:"
            echo ""
            echo "    make -C opam"
            echo ""
            exit 1
        fi
    fi
    rm /tmp/files
done
