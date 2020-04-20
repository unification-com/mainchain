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
                    "/about-mainchain",
                    "/denomination",
                    "/fees-and-gas",
                    "/third-party",
                    "/genesis-settings"
                ]
            },
            {
                title: "Install and Use the Software",
                children: [
                    "/installation",
                    "/accounts-wallets",
                    "/und-commands",
                    "/undcli-commands",
                    "/run-und-as-service"
                ]
            },
            {
                title: "Networks",
                children: [
                  {
                      title: "Mainchain Public TestNet",
                      children: [
                          "/networks/join-testnet",
                          "/networks/become-testnet-validator"
                      ]
                  },
                  {
                      title: "Mainchain Public MainNet",
                      children: [
                          "/networks/join-mainnet",
                          "/networks/become-mainnet-validator"
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
                    "/examples/enterprise-und",
                    "/examples/finchain"
                ]
            },
            {
                title: "In-depth Guides",
                children: [
                    "/guides/cloud/install-aws"
                ]
            }
        ],
    }
}
