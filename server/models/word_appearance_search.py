from .template import SearchModel
from . import util
import pandas as pd 
import numpy as np


class WordAppearance(SearchModel):
    def __init__(self, appearance_min:int=1) -> None:
        super().__init__()
        self.appearance_min = appearance_min


    def train(self) -> None:
        self.files, self.contents = util.load_data()
        self.file_idx = [i for i in range(len(self.files))]
        self._model, self.word_labels = transform_docs_to_mat(self.contents)
    
    def query(self, q: str) -> list:
        _q = util.clean([q])[0]
        q_vec, _ = util.words_to_vec(_q, self.word_labels)
        results = [] 
        appearance = {i:0 for i in range(len(self.files))}
        for i in q_vec:
            results.append(self._model[:,i])
        
        for row_w in results:
            for i, w_c in enumerate(row_w):
                appearance[i] += w_c
        _filtered = {file_idx:cnt for file_idx, cnt in appearance.items() if cnt >= self.appearance_min}
        results = list(dict(sorted(_filtered.items(), key=lambda item: item[1], reverse=True)).keys())      
        results = [self.files[idx] for idx in results]
        return results


def transform_docs_to_mat(contents):
    labels={"NONE":0}
    seq = []
    for con in contents:
        _,labels = util.words_to_vec(con,labels=labels)
    max_len = len(labels)
    labels={"NONE":0}
    for con in contents:
        v,_ = util.words_to_vec(con,labels=labels)
        row = [0]*max_len
        for wt in v:
            row[wt] = 1
        seq.append(row)
    mat = np.array(seq)
    return mat, labels