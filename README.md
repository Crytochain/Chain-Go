# LBRMask Browser Extension for LBR

## Support

If you're a user seeking support, [leave your feedbacks at our GIT site](https://github.com/LBRChain/LBRMask/issues).

中文安装指南可以参考[开发者博客](https://blog.csdn.net/lyq13573221675/article/details/82380846)

安装视频在(YOUTUBE)(https://www.youtube.com/watch?v=kyBJxpa9sA0&t=111s);

## Introduction

In order to help users and developers access LBR blockchain, we modified the [MetaMask Project](https://metamask.io/) to make it work with LBR blockchain. LBR blockchain JSON-RPC is compatiable with Ethereum WEB3 in many methods but is quite different in Transaction Format. Major changes are as the followings:
- Use LBR-tx to replace the ethereumjs-tx for sign a raw transaction object;
- Use LBR-provider-engine to replace the web3-provider-enginer for sending a signed Transaction to LBR network;
- Use LBR-link to provide outside link with LBR explorer for displaying account info.
- Connect with https://gateway.LBR.io instead of infurno.io to provide online services.

LBRMask is a software for users to manage accounts, for sites to easily propose actions to users, and for users to coherently review actions before approving them. We build on this rapidly evolving set of protocols with the goal of empowering the most people to the greatest degree, and aspire to continuously evolve our offering to pursue that goal.


## Developing Compatible Dapps

If you're a web dapp developer, we welcome you to join us to further develop this tool:

### New Dapp Developers

- We recommend this [Learning Solidity](https://karl.tech/learning-solidity-part-1-deploy-a-contract/) tutorial series by Karl Floersch.
- MetaMask team wrote a gentle introduction on [Developing Dapps with Truffle and MetaMask](https://medium.com/metamask/developing-ethereum-dapps-with-truffle-and-metamask-aa8ad7e363ba).

### Current Dapp Developers

- If you have a Dapp on Ethereum, and you want to move to LBR network, you can checkout our [wiki website](https://github.com/LBRChain/LBR-core/wiki/MoveToLBR) for more information. 
- At this moment, LBRMask only supports MotherChain Dapps, MicroChain supports is under developing.

## Building locally

 - Install [Node.js](https://nodejs.org/en/) version 6.3.1 or later.
 - Install dependencies:
   - For node versions up to and including 9, install local dependencies with `npm install`.
   - For node versions 10 and later, install [Yarn](https://yarnpkg.com/lang/en/docs/install/) and use `yarn install`.
 - Install gulp globally with `npm install -g gulp-cli`.
 - Build the project to the `./dist/` folder with `gulp build`.
 - Optionally, to rebuild on file changes, run `gulp dev`.
 - To package .zip files for distribution, run `gulp zip`, or run the full build & zip with `gulp dist`.

 Uncompressed builds can be found in `/dist`, compressed builds can be found in `/builds` once they're built.

### Running Tests

Requires `mocha` installed. Run `npm install -g mocha`.

Then just run `npm test`.

You can also test with a continuously watching process, via `npm run watch`.

You can run the linter by itself with `gulp lint`.

## Development

```bash
npm install
npm start
```

## Build for Publishing

```bash
npm run dist
```

#### Writing Browser Tests

To write tests that will be run in the browser using QUnit, add your test files to `test/integration/lib`.

## Other Docs

- [How to add custom build to Chrome](./docs/add-to-chrome.md)
- [How to add custom build to Firefox](./docs/add-to-firefox.md)
- [How to add new networks to the Provider Menu](./docs/adding-new-networks.md)
- [How to add a new translation to LBRMask](./docs/translating-guide.md)
- [How to develop an in-browser mocked UI](./docs/ui-mock-mode.md)
- [How to develop a live-reloading UI](./docs/ui-dev-mode.md)
- [How to live reload on local dependency changes](./docs/developing-on-deps.md)


