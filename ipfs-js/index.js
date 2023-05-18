import { create } from 'ipfs';
import * as ipns from 'ipns'

async function main() {

    const ipfs = await create();
    const publicKey = "k51qzi5uqu5diel2vj3pzz1g5wbat1qmilgj0z8qb98vwsybplo5ulwwmrtnji"
    const ipnsEntry ="QmcaiYMh28qQCCSesrZ28mkuR2DhVjsYfopJG3jfnUMrRG"

    const topic = 'ipns'
    const receiveMsg = (msg) => console.log(new TextDecoder().decode(msg.data))

    await ipfs.pubsub.subscribe(topic, receiveMsg)
    console.log(`subscribed to ${topic}`)
    }

main();
