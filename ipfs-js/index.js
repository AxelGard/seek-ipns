import * as ipfs from 'ipfs'
import * as ipns from 'ipns'
import { CID } from 'multiformats/cid'


async function main() {

    const ipfs = await ipfs.create();
    const publicKey = "k51qzi5uqu5diel2vj3pzz1g5wbat1qmilgj0z8qb98vwsybplo5ulwwmrtnji"
    const ipnsEntry ="QmcaiYMh28qQCCSesrZ28mkuR2DhVjsYfopJG3jfnUMrRG"

    const topic = '/record/L2lwbnMvEiCGC1J-0c8fai1qZlZ8I5fg8BYN36Tn6tPXsodDl3PTig'
    const receiveMsg = (msg) => console.log(new TextDecoder().decode(msg.data))

    await ipfs.pubsub.subscribe(topic, receiveMsg)
    console.log(`subscribed to ${topic}`)
    }


async function findProvider(){
    const providers = ipfs.dht.findProvider(CID.parse('bafybeibeeinipd62m43ejpp26cn3arkqarao3ddwsot5vglqbpjq2gyb2q'))

    for await (const provider of providers) {
    console.log(provider.id.toString())
    }
}

findProvider();