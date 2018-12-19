from flask import Flask, jsonify, render_template, request, url_for
from dataEngineering import app

@app.route('/')
@app.route('/index')
def index():
   return render_template('index.html' , title='資料工程')

@app.route('/search/<search_text>')
def search(search_text):
   print (search_text)
   return render_template('test.html')