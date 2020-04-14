# Building the Mainchain Documentation

The Mainchain docs are designed to be compiled and distributed using [vuepress](https://vuepress.vuejs.org/).

## Build an deploy

Vuepress can be installed using either `yarn`

```bash
yarn global add vuepress
```

or `npm`

```bash
npm install -g vuepress
```

### Developing

While writing docs, the `vuepress` development server can be used to view changes:

```bash
vuepress dev
```

### Deploying

Run the [deploy-docs.sh](deploy-docs.sh) script in this directory to build and publish changes to https://unification-com.github.io/mainchain:

```bash
./deploy-docs.sh
```


## Format

Docs should be written using standard Markdown, with the exception of a couple of Vuepress specific tags:

```
::: tip
Can be used to render tips
:::
```

```
::: warning
Can be used to render warning notes
:::
```

```
::: danger
Can be used to render danger warning notes
:::
```

These will be rendered and styled by Vuepress in the final html output.

## Sidebar menu

The generate sidebar menu can be edited in [.vurpress/config.js](.vuepress/config.js)

New pages and sections can be added to `themeConfig.sidebar` using the format:

```json
{
    title: "Section Title",
    children: [
        "/some-page",
        "/abother-page"
    ]
}
```

Note the file extension (`.md`, `.html`) is not required.
