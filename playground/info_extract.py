import pandas as pd 
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.cluster import KMeans

# followed a tutorial on TF-IDF https://youtu.be/i74DVqMsRWY

def clean(content):
    ascii_char = [chr(i) for i in range(0,255)]
    numbers = "0123456789"
    non_acc_char =  "\n,.()[]`/:-_*=\\<>|&%@?!\"\'#" + numbers
    non_acc_tokens = ["https","www", "com", "org", "license"]
    for i, _ in enumerate(content):
        for c in non_acc_char:
            content[i] = content[i].replace(c, " ")
        content[i] = content[i].split(" ")
        content[i] = list(filter(lambda c: c != "", content[i]))
        content[i] = [t for t in content[i] if not t in non_acc_tokens ] 
        content[i] = [s.lower() for s in content[i] if all(c in ascii_char for c in s)]

    return content

def TfIdf(content):
    vectorizer = TfidfVectorizer(
        lowercase=True,
        max_features=100,
        max_df=0.8,
        min_df=5,
        ngram_range=(1,3),
        stop_words="english"
    )
    vecs = vectorizer.fit_transform(content)
    feature_names = vectorizer.get_feature_names_out()
    dense = vecs.todense()
    dense_list = dense.tolist()
    all_keywords = []

    for desc in dense_list:
        x=0
        keywords = []
        for word in desc:
            if word > 0:
                keywords.append(feature_names[x]) 
            x += 1
        all_keywords.append(keywords)

    k = 20
    model = KMeans(n_clusters=k, init="k-means++", max_iter=100, n_init=1)
    model.fit(vecs)
    order_ctroides  = model.cluster_centers_.argsort()[:,::-1]
    terms = vectorizer.get_feature_names_out()




def main():
    content = []
    files = ["./cira_README.md", "./pytorch_README.md"]
    for file_path in files: 
        with open(file_path, "r") as f:
            content.append(f.read())


    content = clean(content)
    key_words = []
    for con in content:
        key_words.append(TfIdf(con))



if __name__ == "__main__":
    main()
