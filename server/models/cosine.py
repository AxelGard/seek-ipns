from .template import SearchModel
from . import util
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np

class CosineSimilarity(SearchModel):
    def __init__(self, prob_cut_of:float=0.1, cleaning:bool=True) -> None:
        super().__init__()
        self.prob_cut_of = prob_cut_of
        self.cleaning = cleaning

        self.vectorizer = TfidfVectorizer()
        self.files = None
        self.model = None

    def train(self) -> None:
        self.files, contents = util.load_data()
        if self.cleaning:
            contents = util.clean(contents)
        self.model = self.vectorizer.fit_transform(contents)
    
    def query_proba(self, q: str) -> list:
        assert self.model != None, "train needs to have run before query"
        if self.cleaning:
            q = util.clean([q])[0]
        q_vec = self.vectorizer.transform([q])
        result = [0.0] * len(self.files)
        simi = list(cosine_similarity(q_vec, self.model).tolist()[0])
        _sum = 0.0
        for i,s in enumerate(simi):
            result[i] += s
            _sum += s
        if _sum == 0.0:
            _sum = 1 
        result = [s / _sum for s in result]  
        return result 

    def __str__(self) -> str:
        return "CosineSimilarity"

    def query(self, q: str) -> list: 
        result = []
        p_result = self.query_proba(q)
        idxs = np.array(p_result).argsort()
        for i in idxs:
            if p_result[i] > self.prob_cut_of:
                result.append(self.files[i])
        return result
