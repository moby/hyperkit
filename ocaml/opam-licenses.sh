#!/usr/bin/env sh

set -ex

REPO_ROOT=$(git rev-parse --show-toplevel)
SOURCES=${REPO_ROOT}/_build/sources

mkdir -p $SOURCES

# Collect the license metadata in a central place
OUTPUT=${REPO_ROOT}/ocaml/licenses/

if [ $# = 0 ]; then
  echo "Usage:"
  echo "$0 <package0> ... <packageN>"
  echo "  -- look up license terms for a list of packages"
fi

for DEP in "$@"; do
  echo "Calculating transitive dependencies required by $DEP"
  opam list --required-by "$DEP" --recursive | tail -n +3 > "dependency.$DEP.raw"
  grep -v '\-\-' "dependency.$DEP.raw" | awk '{print $1"."$2}' > "dependency.$DEP"
  rm "dependency.$DEP.raw"
done
cat dependency.* | sort | uniq > all-packages.txt
rm -f "dependency.*"

rm -f "*.extracted"
while read -r PACKAGE; do
  echo "looking for license for ${PACKAGE}"
  URL=$(opam info --field url.src: "$PACKAGE")
  rm -f "/tmp/opam.out"
  if [ -z "${URL}" ]; then
    echo "$PACKAGE has no source: skipping"
    continue
  fi
  echo looking for "$OUTPUT/LICENSE.$PACKAGE.skip"
  if [ -e "$OUTPUT/LICENSE.$PACKAGE.skip" ]; then
    echo "$PACKAGE is not linked: skipping"
    continue
  fi
  if [[ $PACKAGE == *"ppx"* ]]; then
    echo "$PACKAGE is a build dependency: skipping"
    continue
  fi
  if [[ $PACKAGE == *"xen"* ]]; then
    echo "$PACKAGE is not used: skipping"
    continue
  fi
  if [[ $PACKAGE == *"solo5"* ]]; then
    echo "$PACKAGE is not used: skipping"
    continue
  fi
  # if a file has a manually-edited override, use it
  if [ ! -e "$OUTPUT/LICENSE.$PACKAGE" ]; then
    rm -f "$PACKAGE.files"
    DIR=$SOURCES/$PACKAGE
    if [ ! -e "$DIR" ]; then
      ( cd $SOURCES && opam source $PACKAGE )
    fi
    echo "$DIR"
    ls "$DIR" | grep LICENSE >> "$PACKAGE.files" || true
    ls "$DIR" | grep COPYING >> "$PACKAGE.files" || true
    if [ -z "$(cat "$PACKAGE.files")" ]; then
      echo "No LICENSE or COPYING file found in $DIR;"
      echo "please write LICENSE.$PACKAGE yourself"
      exit 1
    fi
  fi
done < all-packages.txt
