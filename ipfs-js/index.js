import * as IPFS from 'ipfs'
import * as ipns from 'ipns'
import { CID } from 'multiformats/cid'


async function main() {

    const ipfs = await IPFS.create();
    const publicKey = "k51qzi5uqu5diel2vj3pzz1g5wbat1qmilgj0z8qb98vwsybplo5ulwwmrtnji"
    const ipnsEntry ="QmcaiYMh28qQCCSesrZ28mkuR2DhVjsYfopJG3jfnUMrRG"
    const myPeerID = "12D3KooWBA3FLioUQPqtj3RT4fxbquGNyb2hfQwXq8UTt5xmxuCi"

    const providers = ipfs.dht.findProvs(CID.parse('QmUVTKsrYJpaxUT7dr9FpKq6AoKHhEM7eG1ZHGL56haKLG'))

    for await (const provider of providers) {
        console.log(provider.id.toString())
    }
}

main()
