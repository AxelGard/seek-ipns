from os import listdir
from os.path import isfile, join
from html.parser import HTMLParser
from bs4 import BeautifulSoup
from . import config

def load_data():
    PATH = config.DATA_STORE
    contents = []
    files = [f"{PATH}{f}" for f in listdir(PATH) if isfile(join(PATH, f))]
    for file_path in files:
        with open(file_path, "r") as f:
            contents.append(f.read())
    return files, contents


def load_stop_words(path: str = "../stopwords.txt"):
    words = []
    with open(path, "r") as f:
        words = f.read().split("\n")
    assert len(words) != 0, "no stop words were found"
    return words


def clean_html(contents:list) -> list:
    result = []
    for html in contents:
        soup = BeautifulSoup(html, 'html.parser')
        result.append(soup.get_text())
    return result

def clean_md(contents:list) -> list: 
    result = []
    for doc in contents: 
        filtered = ""
        for word in doc.split(" "): 
            if not "http" in word:
                filtered += word + " "
        result.append(filtered)
    return result

def clean(contents: list, remove_stop_words=True):
    if config.PARSE_HTML:
        contents = clean_html(contents)
    ascii_char = [chr(i) for i in range(0, 255)]
    numbers = "0123456789"
    non_acc_char = "\n,.()[]{}`/:-_*=\\<>|&%@?!\"'#" + numbers
    non_acc_tokens = ["https", "www", "com", "org", "license"]
    stop_words = [""] 
    if remove_stop_words: 
        stop_words = load_stop_words()
    for i, _ in enumerate(contents):
        for c in non_acc_char:
            contents[i] = contents[i].replace(c, " ")
        contents[i] = contents[i].split(" ")
        contents[i] = list(filter(lambda c: c != "", contents[i]))
        contents[i] = [t for t in contents[i] if not t in non_acc_tokens]
        contents[i] = [
            s.lower() for s in contents[i] if all(c in ascii_char for c in s)
        ]
        contents[i] = [t for t in contents[i] if not t in stop_words]
        for j, word in enumerate(contents[i]):
            if word[-1] == "s":
                contents[i][j] = word[:-1]
    return [" ".join(con) for con in contents]


def words_to_vec(words: str, labels: dict = {}):
    _words = words.split(" ")
    vec = []
    for word in _words:
        if word not in labels:
            labels[word] = len(labels)
        vec.append(labels[word])
    return vec, labels
