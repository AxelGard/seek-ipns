import pandas as pd 
from sklearn.feature_extraction.text import TfidfVectorizer
from os import listdir
from os.path import isfile, join
from collections import defaultdict

PATH = "../data/"
contents = []
files = [f"{PATH}{f}" for f in listdir(PATH) if isfile(join(PATH, f))]
tf_idf_data = None 
data = {}

def model_init(): 
    global contents 
    global files
    global tf_idf_data
    global data 

    for file_path in files: 
        with open(file_path, "r") as f:
            contents.append(f.read())
    
    contents = clean(contents)
    tf_idf_data = TfIdf(contents)
    data = tf_idf_data.to_dict()
    print("DONE WITH TRAINING")
    


def load_stop_words(path:str="../stopwords.txt"):
    words = []
    with open(path, "r") as f: 
       words = f.read().split("\n")
    assert len(words) != 0, "no stop words were found"
    return words


def clean(content:list):
    ascii_char = [chr(i) for i in range(0,255)]
    numbers = "0123456789"
    non_acc_char =  "\n,.()[]{}`/:-_*=\\<>|&%@?!\"\'#" + numbers
    non_acc_tokens = ["https","www", "com", "org", "license"]
    stop_words = load_stop_words()
    for i, _ in enumerate(content):
        for c in non_acc_char:
            content[i] = content[i].replace(c, " ")
        content[i] = content[i].split(" ")
        content[i] = list(filter(lambda c: c != "", content[i]))
        content[i] = [t for t in content[i] if not t in non_acc_tokens ] 
        content[i] = [s.lower() for s in content[i] if all(c in ascii_char for c in s)]
        content[i] = [t for t in content[i] if not t in stop_words] 
        for j, word in enumerate(content[i]):
            if word[-1] == "s": 
                content[i][j] = word[:-1]
    return [" ".join(con) for con in content]


def TfIdf(content:list):
    vectorizer = TfidfVectorizer()
    vecs = vectorizer.fit_transform(content)
    feature_names = vectorizer.get_feature_names_out()
    dense = vecs.todense()
    dense_list = dense.tolist()
    df = pd.DataFrame(dense_list, columns=feature_names)
    return df

def query_data(tf_idf_data:dict, query:str)->list:
    query = clean([query])[0]
    words_tf_idf = {}
    query_words = query.split(" ")
    for word in query_words: 
        if word in tf_idf_data.keys():
            words_tf_idf[word] = tf_idf_data[word]

    result = []
    highest = 0.0
    for word, tf_idf in words_tf_idf.items():
        for idx, val in tf_idf.items(): 
            if val > highest: 
                highest = val 
                result.append((val, idx))
    result.sort()
    result.reverse()
    lookup = set()
    result = [idx for _, idx in result if idx not in lookup and lookup.add(idx) is None]
    return result



def search_query(q:str)->list:
    result = []
    global data 
    global files
    idxs = query_data(data, q)
    for idx in idxs: 
        file = files[idx]
        file = file[2:]
        result.append(file)
    return result
