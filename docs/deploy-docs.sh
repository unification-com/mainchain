#!/usr/bin/env sh

# abort on errors
set -e

rm -rf .vuepress/dist

vuepress build

echo "docs.unification.io" > .vuepress/dist/CNAME

cd .vuepress/dist

git init

git add -A

git commit -m 'deploy'

git push -f git@github.com:unification-com/mainchain.git master:gh-pages
