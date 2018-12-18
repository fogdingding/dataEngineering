from flask import Flask, jsonify, render_template, request, url_for
from app import app

@app.route('/')
@app.route('/index')
def index():
    subject = {'subject': '資料工程'}
    return render_template('index.html' , title='資料工程', subject=subject)

@app.route('/search')
def add_numbers():
   pass