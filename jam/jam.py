from jam import app
from os import path
import json

from flask import render_template, request, redirect, url_for

BASE_PATH = path.dirname(path.realpath(__file__))
STATE_DEFAULT_PATH = path.join(BASE_PATH, "static", "state_default.json")
STATE_PATH = path.join(BASE_PATH, "static", "state.json")

@app.route('/', methods=['GET'])
def index():
    with open(STATE_PATH, 'r') as fin:
        data = json.load(fin)
        started = data['gameStarted']
        cards = data['cards']
    
    if not started:
        return render_template("welcome.html")
    return render_template("home.html", cards=cards)

@app.route('/start_game', methods=['POST'])
def start_game():
    name = request.form['name']
    with open(STATE_PATH, 'r') as fin:
        data = json.load(fin)
        data['gameStarted'] = True
        data['player1']['name'] = name
    with open(STATE_PATH, 'w') as fout:
        json.dump(data, fout)
    return redirect(url_for('index'))

@app.route('/reset_game', methods=['POST'])
def reset_game():
    with open(STATE_DEFAULT_PATH, 'r') as fin:
        data = json.load(fin)
    with open(STATE_PATH, 'w') as fout:
        json.dump(data, fout)
    return redirect(url_for('index'))