from flask import Flask
from .models import RandomForest as Model
import json
from .models import score

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

@app.get("/api/model/score")
def get_score():
    scores = score.main()
    return json.dumps(scores)
