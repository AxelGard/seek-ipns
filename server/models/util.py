from os import listdir
from os.path import isfile, join

def load_data():
    PATH = "../data/"
    contents = []
    files = [f"{PATH}{f}" for f in listdir(PATH) if isfile(join(PATH, f))]
    for file_path in files: 
        with open(file_path, "r") as f:
            contents.append(f.read())
    return files, contents


def load_stop_words(path:str="../stopwords.txt"):
    words = []
    with open(path, "r") as f: 
       words = f.read().split("\n")
    assert len(words) != 0, "no stop words were found"
    return words


def clean(contents:list):
    ascii_char = [chr(i) for i in range(0,255)]
    numbers = "0123456789"
    non_acc_char =  "\n,.()[]{}`/:-_*=\\<>|&%@?!\"\'#" + numbers
    non_acc_tokens = ["https","www", "com", "org", "license"]
    stop_words = load_stop_words()
    for i, _ in enumerate(contents):
        for c in non_acc_char:
            contents[i] = contents[i].replace(c, " ")
        contents[i] = contents[i].split(" ")
        contents[i] = list(filter(lambda c: c != "", contents[i]))
        contents[i] = [t for t in contents[i] if not t in non_acc_tokens ] 
        contents[i] = [s.lower() for s in contents[i] if all(c in ascii_char for c in s)]
        contents[i] = [t for t in contents[i] if not t in stop_words] 
        for j, word in enumerate(contents[i]):
            if word[-1] == "s": 
                contents[i][j] = word[:-1]
    return [" ".join(con) for con in contents]


def words_to_vec(words:str, labels:dict={}):
    _words = words.split(" ")
    vec = []
    for word in _words:
        if word not in labels:
            labels[word]=len(labels)
        vec.append(labels[word])
    return vec, labels