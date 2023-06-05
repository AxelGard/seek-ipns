from .tfidf_search import TfIdf
from .cosine import CosineSimilarity
from .svm import SupportVectorMachine
from .random_forest import RandomForest

ALL = [TfIdf, CosineSimilarity, SupportVectorMachine, RandomForest]


from sklearn.metrics import f1_score
from sklearn.metrics import accuracy_score
from sklearn.metrics import average_precision_score
from sklearn.metrics import ndcg_score
from . import util
from . import validation_qa
import numpy as np
import copy


query_results_answer = {}

FILES, _ = util.load_data()


def get_f1_scores(model):
    global query_results_answer
    scores = {}
    for query, expec_files in query_results_answer.items():
        results = model.query(query)
        results = files_to_vec(results)
        f1 = f1_score(expec_files, results)
        scores[query] = f1
    return scores


def get_accuracy_scores(model):
    global query_results_answer
    scores = {}
    for query, expec_files in query_results_answer.items():
        results = model.query(query)
        results = files_to_vec(results)
        acc = accuracy_score(expec_files, results)
        scores[query] = acc
    return scores


def get_average_precision_scores(model):
    global query_results_answer
    scores = {}
    for query, expec_files in query_results_answer.items():
        results = model.query_proba(query)
        acc = average_precision_score(expec_files, results)
        scores[query] = acc
    return scores

def get_ndcg_scores(model):
    global query_results_answer
    scores = {}
    for query, expec_files in query_results_answer.items():
        results = model.query_proba(query)
        ndcg = ndcg_score([expec_files], [results])
        scores[query] = ndcg
    return scores
    


def files_to_vec(files):
    ans_q = [0] * len(FILES)
    for f in files:
        if f[:2] != "..":
            f = ".." + f
        ans_q[FILES.index(f)] = 1
    return ans_q


def get_scores(model):
    scores = {}
    model.train()
    f1_scores = get_f1_scores(model)
    acc_scores = get_accuracy_scores(model)
    prec_scores = get_average_precision_scores(model)
    ndcg_scores = get_ndcg_scores(model)
    scores = {
        "ndcg_mean": np.array(list(ndcg_scores.values())).mean(),
        #"prec_scores":prec_scores,
        "prec_mean": np.array(list(prec_scores.values())).mean(),
        #"prec_scores":prec_scores,
        "f1_mean": np.array(list(f1_scores.values())).mean(),
        # "f1_scores":f1_scores,
        "acc_mean": np.array(list(acc_scores.values())).mean(),
        #"acc_scores":acc_scores,
    }
    return scores

def setup():
    global query_results_answer
    query_results_answer = copy.deepcopy(validation_qa.GITHUB_README)
    for query, expec_files in query_results_answer.items():
        ans_q = files_to_vec(expec_files)
        query_results_answer[query] = ans_q


def main():
    setup()
    scores = {}
    for Model in ALL:
        model = Model()
        model.train()
        scores[str(model)] = get_scores(model)
    return scores
    # scores = sorted(scores, key=lambda x: x["prec_mean"], reverse=True)