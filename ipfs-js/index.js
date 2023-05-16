import { create } from 'ipfs';
import all from 'it-all'

async function getLsCid() {
    const ipfs = await create();
    const cid = "bafybeiarvsvlo6hmb6v5uiz6zm2c7hssv6smctoiddya6bemixagmsbcdu"
    for await (const file of ipfs.ls(cid)) {
      console.log(file)
    }
}

async function main() {
    const ipfs = await create();
    const peerId = "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"
    const list = await ipfs.bitswap.wantlistForPeer(peerId)
    console.log(list)
}

main();
