import pandas as pd 
from sklearn.feature_extraction.text import TfidfVectorizer
from .template import SearchModel
from . import util

class TfIdf(SearchModel):
    def __init__(self, tf_idf_limit:float=0.01) -> None:
        super().__init__()
        self.vectorizer = TfidfVectorizer()
        self.features = None 
        self.files, self.contents = util.load_data()
        self.tf_idf = None
        self.tf_idf_limit = tf_idf_limit

    def train(self)-> None:
        self.contents = util.clean(self.contents)
        self.features = self.vectorizer.fit_transform(self.contents)
        self.feature_names = self.vectorizer.get_feature_names_out()
        dense = self.features.todense()
        dense_list = dense.tolist()
        self.tf_idf = pd.DataFrame(dense_list, columns=self.feature_names).to_dict()

    def _query(self, q:str) -> list: 
        assert self.tf_idf != None, "the model needs to be trained before running queries "
        _q = util.clean([q])[0]
        words_tf_idf = {}
        query_words = _q.split(" ")
        for word in query_words: 
            if word in self.tf_idf.keys():
                words_tf_idf[word] = self.tf_idf[word]

        result = []
        for word, tf_idf in words_tf_idf.items():
            for idx, val in tf_idf.items(): 
                if val >= self.tf_idf_limit: 
                    result.append((val, idx))
        result.sort()
        result.reverse()
        lookup = set()
        result = [idx for _, idx in result if idx not in lookup and lookup.add(idx) is None]
        return result

    def query(self, q: str) -> list: 
        result = []
        idxs = self._query(q)
        for idx in idxs: 
            file = self.files[idx]
            file = file[2:]
            result.append(file)
        return result
