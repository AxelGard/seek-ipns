from .template import SearchModel
from . import util
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.svm import SVC

class SupportVectorMachine(SearchModel):
    def __init__(self, kernel:str="linear", C:float=1.0, prob_min:float=0.1005 ,clean:bool=True) -> None:
        self.vectorizer = TfidfVectorizer()
        self.clean=clean
        self.model = None
        self.kernel=kernel
        self.C = C
        self.prob_min = prob_min
        
    def train(self) -> None:
        self.files, contents = util.load_data()
        if self.clean:
            contents = util.clean(contents=contents)
        features = self.vectorizer.fit_transform(contents)
        self.model = SVC(probability=True, kernel=self.kernel, C=self.C)
        self.model.fit(features, self.files)

    def query_proba(self, q: str) -> list:
        assert self.model != None, "train needs to have run before query"
        if self.clean:
            q = util.clean([q])[0]
        q_vec = self.vectorizer.transform([q])
        prob_of_file = self.model.predict_proba(q_vec)[0]
        return prob_of_file

    def query(self, q: str) -> list:
        prob_of_file = self.query_proba(q)
        result = []
        result = [(prob, self.files[idx]) for idx, prob in enumerate(prob_of_file) if prob > self.prob_min]
        result.sort()
        result.reverse()
        result = [f for _, f in result ]
        return result
    
