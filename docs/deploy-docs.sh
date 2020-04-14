#!/usr/bin/env sh

# abort on errors
set -e

rm -rf .vuepress/dist

vuepress build

cd .vuepress/dist

git init

git add -A

git commit -m 'deploy'

git push -f git@github.com:unification-com/mainchain.git master:gh-pages
