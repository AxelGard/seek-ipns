import { create } from 'ipfs';

async function main() {
  try {
    const ipfs = await create();
    const peers = await ipfs.pubsub.peers();
    const targetPeerId = 'QmcMdoQ3vDUpd1uXc3XUyzDcF6e214pUii5wXQBL5nV2xq';
    const targetPeers = peers.filter(peer => peer.includes(targetPeerId));

    for (const peer of targetPeers) {
      const cidList = await ipfs.ls(peer);
      console.log(`CIDs for peer ${peer}:`, cidList);
    }
  } catch (error) {
    console.error('An error occurred:', error);
  }
  console.log("DONE")
}

main();
