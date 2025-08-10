import random

from containers import routes as container_routes
from images import routes as image_routes

from flask import Flask, render_template

app = Flask(__name__)

@app.route("/")
def index():
    return render_template("index.html")

@app.route("/more-items")
def more_items():
    items = []
    for _ in range(0, random.randint(1, 11)):
        items.append("".join(random.choices("abcdefghijklmnopqrstuvwxyz", k=random.randint(4, 11))))
    return render_template("partials/more_items.html", items=items)


if __name__ == "__main__":
    app.run(debug=True)