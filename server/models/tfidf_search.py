import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from .template import SearchModel
from . import util
import numpy as np
from .config import USE_METADATA
from .metadata import apply_metadata


class TfIdf(SearchModel):
    def __init__(self, prob_cut_of: float = 0.1, clean: bool = True) -> None:
        super().__init__()
        self.vectorizer = TfidfVectorizer()
        self.features = None
        self.files, self.contents = util.load_data()
        self.tf_idf = None
        self.prob_cut_of = prob_cut_of
        self.clean = clean

    def train(self) -> None:
        self.contents = util.clean(self.contents)
        self.features = self.vectorizer.fit_transform(self.contents)
        self.feature_names = self.vectorizer.get_feature_names_out()
        dense = self.features.todense()
        self.model = dense.tolist()
        self.tf_idf = pd.DataFrame(self.model, columns=self.feature_names).to_dict()

    def query_proba(self, q: str) -> list:
        assert (
            self.tf_idf != None
        ), "the model needs to be trained before running queries "
        if self.clean:
            q = util.clean([q])[0]
        words_tf_idf = {}
        query_words = q.split(" ")
        for word in query_words:
            if word in self.tf_idf.keys():
                words_tf_idf[word] = self.tf_idf[word]
        result = [0] * len(self.files)
        _sum = 0.0
        for word, tf_idf in words_tf_idf.items():
            for idx, val in tf_idf.items():
                result[idx] += val
                _sum += val
        if _sum == 0.0:
            _sum = 1
        result = [v / _sum for v in result]
        return result

    def query(self, q: str) -> list:
        result = []
        p_result = self.query_proba(q)
        intermediate_results = []
        idxs = np.array(p_result).argsort()
        for i in idxs:
            if p_result[i] > self.prob_cut_of:
                intermediate_results.append([p_result[i],self.files[i]])
        intermediate_results.sort()
        intermediate_results.reverse()
        if USE_METADATA: 
            intermediate_results = apply_metadata(intermediate_results)
        result = [f for _, f in intermediate_results]
        return result

    def __str__(self) -> str:
        return "TfIdf"
