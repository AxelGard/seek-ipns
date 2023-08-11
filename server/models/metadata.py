import math 
import pandas as pd
import numpy as np
from . import config
import pprint

def number_of_peers(cid:str) -> int: 
    df = pd.read_csv(config.META_DATA_STORE +"week1/number_of_hosts.csv") 
    row = df[df["cid"] == cid]
    if len(row) == 1:
        peers = row["other_peers"]
        if peers.hasnans: return 0
        peers = peers.to_list()[0] 
        return len([p for p in peers.split(",")])
    else: # should not happened but you never know
        _sum = 0
        peers = row["other_peers"].to_list()
        for pp in peers:
            _sum += len([p for p in pp.split(",")])
        return _sum
    

def life_time(cid:str) -> int: 
    _time = 0
    for i in range(0,4):
        df = pd.read_csv(config.META_DATA_STORE +f"week{4-i}/time.csv") 
        row = df[df["cid"] == cid]
        if len(row) == 0:
            continue
        else: 
            _time += 1
    print("time", _time)
    return _time


def size(cid:str) -> int: 
    pass


def sigmoid(x:int) -> float:
    return 1 / (1 + (math.e**-x))

def inverse(x:int) -> float:
    return 1 / (1 + x**2)

def apply_metadata(ranking:list) -> list:
    pp = pprint.PrettyPrinter(indent=4)
    pp.pprint(ranking)
    new_ranking = []
    for p,f in ranking: 
        cid = f.split("/")[-1]
        p *= inverse(number_of_peers(cid))
        p *= sigmoid(life_time(cid))
        new_ranking.append([p,f])
    new_ranking.sort()
    new_ranking.reverse()
    print("")
    pp.pprint(new_ranking)
    return new_ranking


#print(number_of_peers("QmaiqBUNA1SfS1rShgAHbrJCTQBgwKRKkKQiNBc4qECmdR"))