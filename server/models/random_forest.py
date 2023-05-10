from .template import SearchModel
from . import util
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.ensemble import RandomForestClassifier

class RandomForest(SearchModel):
    def __init__(self, clean:bool=True, n_estimators:int=1000, prob_min:float=0.05) -> None:
        self.clean = clean
        self.prob_min = prob_min
        self.vectorizer = TfidfVectorizer()
        self.model = RandomForestClassifier(n_estimators=n_estimators)

    def train(self) -> None:
        self.files, contents = util.load_data()
        if self.clean:
            contents = util.clean(contents)
        self.files_idxs = [i for _, i in enumerate(self.files)]
        self.features  = self.vectorizer.fit_transform(contents)
        self.model.fit(self.features, self.files_idxs)
    
    def query(self, q: str) -> list:
        if self.clean:
            q = util.clean([q])[0]
        q_vec = self.vectorizer.transform([q])
        prob_of_file = self.model.predict_proba(q_vec)[0]
        result = [(prob,self.files[idx]) for idx, prob in enumerate(prob_of_file) if  prob > self.prob_min]
        result.sort()
        result.reverse()
        result = [f for _, f in result]
        return result


