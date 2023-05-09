
from .tfidf_search import TfIdf
from .word_appearance_search import WordAppearance
from .cosine import CosineSimilarity
from .svm import SupportVectorMachine

ALL = [TfIdf, SupportVectorMachine]


from sklearn.metrics import f1_score
from . import util



query_results_answer = {
    "algorithmic trading" : [
        '../data/cira_README.md',
    ],
    "JavaScript library for building user interfaces" : [
        '../data/svelte_README.md',
        "../data/react_README.md",
        "../data/vuejs_README.md",
    ],
}

FILES, _ = util.load_data()

def get_scores(Model):
    scores = {}
    model = Model()
    model.train()
    for query, expec_files in query_results_answer.items():
        results = model.query(query)
        results = files_to_vec(results)
        f1 = f1_score(expec_files, results)
        scores[query] = f1
    return scores

def files_to_vec(files):
    ans_q = [0]*len(FILES)
    for f in files: 
        if f[:2] != "..":
            f = ".." + f
        ans_q[FILES.index(f)] = 1
    return ans_q

def run_score_test():
    scores = []
    for model in ALL:
        scores.append({
            "model":str(model), 
            "scores": get_scores(model)})
    return scores

def main():
    for query, expec_files in query_results_answer.items():
        ans_q = files_to_vec(expec_files)
        query_results_answer[query] = ans_q
    scores = run_score_test()
    return scores