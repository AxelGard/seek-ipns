
class SearchModel:
    def __init__(self) -> None:
        """ setup vectorizer and other needed classes and pre train functions """
        pass 

    def train(self):
        """ fit and train model on data so that it is ready to run """
        raise NotImplementedError

    def query(self, q:str)->list:
        """ takes a search query and returns a list of files,
        where the first file in the list is the most correlated
        result to the given query """
        raise NotImplementedError