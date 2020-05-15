module.exports = {
    title: 'Mainchain Documentation',
    description: 'Welcome to the documentation for Unification Mainchain',
    base: '/',
    markdown: {
        // options for markdown-it-toc
        toc: {includeLevel: [2, 3]}
    },

    themeConfig: {
        lastUpdated: 'Last Updated',
        repo: 'unification-com/mainchain',
        docsDir: 'docs',
        logo: '/assets/img/unification_logoblack.png',
        nav: [{
            text: 'Releases',
            link: 'https://github.com/unification-com/mainchain/releases'
        }],
        sidebar: [
            {
                title: "About Mainchain",
                children: [
                    "/introduction/about-mainchain",
                    "/introduction/denomination",
                    "/introduction/fees-and-gas",
                    "/introduction/total-supply",
                    "/introduction/genesis-settings",
                    "/introduction/delegators",
                    "/introduction/validators",
                    "/introduction/faqs"
                ]
            },
            {
                title: "Install and Use the Software",
                children: [
                    "/software/installation",
                    "/software/accounts-wallets",
                    "/networks/join-network",
                    "/software/run-und-as-service",
                    "/software/light-client-rpc",
                    {
                      title: "CLI Command & Config References",
                      children: [
                        "/software/und-commands",
                        "/software/undcli-commands",
                        "/software/und-mainchain-config-ref",
                        "/software/und-mainchain-app-config-ref"
                      ]
                    }
                ]
            },
            {
                title: "Networks",
                children: [
                  {
                      title: "Mainchain Public Networks",
                      children: [
                          "/networks/join-network",
                          "/networks/become-validator"
                      ]
                  },
                  {
                      title: "Play with DevNet",
                      children: [
                          "/networks/local-devnet"
                      ]
                  },
                  "/networks/participation",
                ]
            },
            {
                title: "Tx & Query Examples",
                children: [
                    "/examples/transactions",
                    "/examples/wrkchain",
                    "/examples/beacon",
                    "/examples/enterprise-fund",
                    "/examples/finchain"
                ]
            },
            {
                title: "In-depth Guides",
                children: [
                    "/guides/cloud/install-aws",
                    "/guides/cloud/install-gc"
                ]
            },
            {
                title: "Developers",
                children: [
                    "/developers/third-party",
                ]
            },
        ],
    }
}
