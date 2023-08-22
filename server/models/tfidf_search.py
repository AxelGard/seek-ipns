import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from .template import SearchModel
from . import util
import numpy as np
from .config import USE_METADATA
from .metadata import apply_metadata


class TfIdf(SearchModel):
    def __init__(self, cut_of: float = 0.1, clean: bool = True) -> None:
        super().__init__()
        self.vectorizer = TfidfVectorizer()
        self.features = None
        self.files, self.contents = util.load_data()
        self.tf_idf = None
        self.cut_of = cut_of
        self.clean = clean

    def train(self) -> None:
        self.contents = util.clean(self.contents)
        self.features = self.vectorizer.fit_transform(self.contents)
        self.feature_names = self.vectorizer.get_feature_names_out()
        dense = self.features.todense()
        self.model = dense.tolist()
        self.tf_idf = pd.DataFrame(self.model, columns=self.feature_names)

    def query_proba(self, q: str) -> dict:
        if self.clean:
            q = util.clean([q])[0]
        result = []
        tf_idf = self.tf_idf.copy()
        q = q.split(" ")
        q = tf_idf.columns[tf_idf.columns.isin(q)]
        tf_idf = tf_idf[q].T
        document_wights = tf_idf.sum(axis=0)
        _sum = document_wights.sum()
        if _sum == 0.0:
            _sum = 1
        result = [v / _sum for v in document_wights.values]
        return result

    def query(self, q: str) -> list:
        result = []
        p_result = self.query_proba(q)
        intermediate_results = []
        idxs = np.array(p_result).argsort()
        for i in idxs:
            if p_result[i] > self.cut_of:
                intermediate_results.append([p_result[i],self.files[i]])
        intermediate_results.sort()
        intermediate_results.reverse()
        if USE_METADATA: 
            intermediate_results = apply_metadata(intermediate_results)
        result = [f for _, f in intermediate_results]
        return result

    def __str__(self) -> str:
        return "TfIdf"
