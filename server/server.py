from flask import Flask
from .models import WordAppearance as Model
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
    for i,f in enumerate(files): 
        result.append({
            "file":f,
            "cid": "",
            "info": "",
            "meta_data":{},
            "rank":i
            })
    return json.dumps(result) 
