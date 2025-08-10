import random

from containers import routes as c_routes
from images import routes as i_routes

from flask import Flask, render_template

app = Flask(__name__)
app.register_blueprint(c_routes.containers_bp)
app.register_blueprint(i_routes.images_bp)

@app.route("/")
def index():
    return render_template("index.html")



if __name__ == "__main__":
    app.run(debug=True)