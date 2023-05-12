
from .tfidf_search import TfIdf
from .cosine import CosineSimilarity
from .svm import SupportVectorMachine
from .random_forest import RandomForest

ALL = [TfIdf, CosineSimilarity, SupportVectorMachine, RandomForest]


from sklearn.metrics import f1_score
from sklearn.metrics import accuracy_score
from sklearn.metrics import average_precision_score 
from . import util
import numpy as np


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

def files_to_vec(files):
    ans_q = [0]*len(FILES)
    for f in files: 
        if f[:2] != "..":
            f = ".." + f
        ans_q[FILES.index(f)] = 1
    return ans_q

def run_score_test():
    scores = {}
    for Model in ALL:
        model = Model()
        model.train()
        f1_scores = get_f1_scores(model)
        acc_scores = get_accuracy_scores(model)
        prec_scores = get_average_precision_scores(model)
        scores[str(model)] = {
            "prec_mean": np.array(list(prec_scores.values())).mean(),
            #"prec_scores":prec_scores,
            "f1_mean": np.array(list(f1_scores.values())).mean(),
            #"f1_scores":f1_scores,
            "acc_mean": np.array(list(acc_scores.values())).mean(),
            #"acc_scores":acc_scores,
        }

    return scores

def main():
    global query_results_answer

    query_results_answer = {
        "algorithmic trading" : [
            '../data/cira_README.md',
        ],
        "a of buying a part of a company with code" : [
            '../data/cira_README.md',
        ],
        "JavaScript library for building user interfaces" : [
            '../data/svelte_README.md',
            "../data/react_README.md",
            "../data/vuejs_README.md",
        ],
        "A memory safe programming language": [ 
            "../data/rust_README.md",
        ],
        "A compiler for my C program":[
            "../data/gcc_README.txt",
        ], 
        "make a low level language to machine code":[
            "../data/gcc_README.txt",
        ], 
        "framework or library for ML and AI": [
            "../data/pytorch_README.md", 
            "../data/tensorflow_README.md",
        ],
        "framework or library for Machine learning and artificial intelligence": [
            "../data/pytorch_README.md", 
            "../data/tensorflow_README.md",
        ],
        "a simple programming language":[
            "../data/cpython_README.rst", 
            "../data/rust_README.md",
        ],
        "I need a operating system":[
            "../data/linux_README.txt",
        ]
    }

    for query, expec_files in query_results_answer.items():
        ans_q = files_to_vec(expec_files)
        query_results_answer[query] = ans_q
    scores = run_score_test()
    #scores = sorted(scores, key=lambda x: x["prec_mean"], reverse=True)
    return scores