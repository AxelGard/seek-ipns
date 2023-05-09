from .template import SearchModel
from . import util
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity


class CosineSimilarity(SearchModel):
    def __init__(self, similarity_threshold:float=0.0, cleaning:bool=True) -> None:
        super().__init__()
        self.similarity_threshold = similarity_threshold
        self.cleaning = cleaning

        self.vectorizer = TfidfVectorizer()
        self.files = None
        self.model = None

    def train(self) -> None:
        self.files, contents = util.load_data()
        if self.cleaning:
            contents = util.clean(contents)
        self.model = self.vectorizer.fit_transform(contents)
  

    def query(self, q: str) -> list:
        assert self.model != None, "train needs to have run before query"
        results = []
        if self.cleaning:
            q = util.clean([q])[0]
        q_vec = self.vectorizer.transform([q])
        simi = cosine_similarity(q_vec, self.model)
        idxs = simi.argsort()[0]
        results = [self.files[idx] for idx in idxs if simi[0][idx] > self.similarity_threshold]
        return results

