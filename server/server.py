from flask import Flask, render_template, request
from .models.tfidf_search import TfIdf as Model
import json

app = Flask(__name__)
model = Model()
model.train()

@app.route("/")
def landing():
    return {"msg":"GitBay API"}


@app.get("/api/query/<q>")
def query(q):
    files = model.query(q)
    result = []
    for f in files: 
        result.append({
            "file":f,
            "cid": "",
            "info": "",
            "meta_data":{}
            })
    return json.dumps(result) 
