import sys
import os

def is_online(): 
    hello_world_test = os.popen("ipfs cat QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o").read()
    assert hello_world_test == "hello world\n", "Could not access IPFS daemon"

def main():
    args = sys.argv[1:]
    key = ""
    if len(args) > 1 and "--key=" in args[1]:
        key = args[1]
    path = args[0]
    if os.path.exists(path):
        added = os.popen(f"ipfs add -r -H {path}").read()
        print(added)
        cid = added.split("added")[-1].split(" ")[1]
        pub_cmd = f"ipfs name publish /ipfs/{cid} " + key
        ipns_pub = os.popen(pub_cmd).read()
        print(ipns_pub)
        return
    raise FileExistsError
        
if __name__ == "__main__":
    is_online()
    main()