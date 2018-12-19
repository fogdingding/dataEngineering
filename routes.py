from flask import Flask, jsonify, render_template, request, url_for
from dataEngineering import app
from subprocess import Popen,PIPE,STDOUT,call

@app.route('/')
# 這邊為主畫面
@app.route('/index')
def index():
   # render_template 內建可引導至templates資料夾的index.html，並傳送title='資料工程'。
   return render_template('index.html' , title='資料工程')

#這邊為搜尋完的頁面
@app.route('/search/<search_text>')
def search(search_text):
   #gerp內建搜尋函數
   search_string = 'grep '+search_text+' .\data\output.tsv -n'
   #這邊使用系統指令，並從中取出系統印出之字串
   proc=Popen(search_string, shell=True, stdout=PIPE, )
   output=proc.communicate()[0]
   #字串處理
   output_list = output.decode('utf-8').split("\n")
   return render_template('search.html',search_text=output_list)