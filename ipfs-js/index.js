import { create } from 'ipfs';
import all from 'it-all'

async function main() {
    const ipfs = await create();
    const cid = "bafybeiarvsvlo6hmb6v5uiz6zm2c7hssv6smctoiddya6bemixagmsbcdu"
    for await (const file of ipfs.ls(cid)) {
      console.log(file)
    }
}

main();
