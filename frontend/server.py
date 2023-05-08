from flask import Flask, render_template, request
import requests

app = Flask(__name__, static_url_path='/templates')

@app.route("/")
def landing():
    return render_template("index.html") 


@app.get("/query")
def query():
    q = request.args.get("q")
    r = requests.get(f"http://127.0.0.1:5000/api/query/{q}")
    if r.status_code != 200:
        return r.status_code
    files = []
    for data in r.json():
        files.append(data["file"])
    return render_template("results.html", files=files, search_query=q)


@app.get("/data/<file>")
def serve_file(file):
    file_con = None 
    with open(f"../data/{file}", "r") as f: 
        file_con = f.read()
    return file_con


if __name__ == '__main__':
    app.run(host="localhost", port=8000, debug=True)