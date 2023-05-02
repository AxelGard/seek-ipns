from flask import Flask, render_template, request
from .tfidf_search import search_query, model_init

app = Flask(__name__, static_url_path='/templates')
model_init()

@app.route("/")
def landing():
    return render_template("index.html") 


@app.get("/query")
def query():
    q = request.args.get("q")
    files = search_query(q)
    return render_template("results.html", files=files)


@app.get("/data/<file>")
def serve_file(file):
    file_con = None 
    with open(f"./data/{file}", "r") as f: 
        file_con = f.read()
    return file_con