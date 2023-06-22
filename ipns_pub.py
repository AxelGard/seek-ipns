import sys
import os

def main():
    args = sys.argv[1:]
    key = ""
    if args[0] == "-h" or args[0] == "--help":
        help()
        return
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
        ipns_recod = ipns_pub.split(":")[0].split(" ")[-1]
        print(f"Sharable link: https://ipfs.io/ipns/{ipns_recod}")
        print("")

        return
    raise FileExistsError
        
def is_online(): 
    hello_world_test = os.popen("ipfs cat QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o").read()
    assert hello_world_test == "hello world\n", "Could not access IPFS daemon"

def help():
    help_msg = """Usage:
    $python3 ipns_pub.py <path to dir or file>

    If you want a other key then the default key then:
    $python3 ipns_pub.py <path to dir or file> --key=<key name>

    This will publish the given path to IPNS
    """
    print(help_msg)

if __name__ == "__main__":
    is_online()
    main()