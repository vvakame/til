#!/bin/sh -eux

cd `dirname $0`

if [[ ! -x `which goapp` ]]; then
  echo "goapp is not int PATH!"
  exit 1
fi

rm -rf src/public
cp -r ../front/dist/ src/public/

gb gae serve src
