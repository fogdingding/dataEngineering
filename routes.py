from flask import Flask, jsonify, render_template, request, url_for
from dataEngineering import app

@app.route('/')
@app.route('/index')
def index():
   return render_template('index.html' , title='資料工程')

@app.route('/search')
def add_numbers():
   pass