from flask import Flask, jsonify, render_template, request, url_for
from dataEngineering import app
from subprocess import Popen,PIPE,STDOUT,call
import base64
import pandas as pd
from pandas_datareader import data as web
import datetime
import matplotlib.pyplot as plt 

# 圖片轉乘base64的寫法。
# import base64
# with open('test.png','rb') as img_f:
#       img_stream = img_f.read()
#       img_stream = base64.b64encode(img_stream)
# print(str(img_stream,'utf-8'))

#這邊為把歷史股價抓下來存成圖片。
def save_image(stock):
   start = datetime.datetime(2018,1,1)
   end = datetime.date.today()
   data = web.DataReader(stock, "yahoo", start,end)
   c = data['Close']
   plt.cla()
   ax=c.plot(title=stock)
   fig = ax.get_figure()
   fig.savefig('./data/'+stock+'.png')

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
   search_string = 'grep '+search_text+' .\data\output.ssv'
   #這邊使用系統指令，並從中取出系統印出之字串
   proc=Popen(search_string, shell=True, stdout=PIPE, )
   output=proc.communicate()[0]
   #字串處理
   output_list = output.decode('utf-8').split("\n")
   return render_template('search.html',search_text=output_list)

#這邊為把圖片轉成base64傳到前端
@app.route('/image/<id>')
def get_image(id):
   img_stream = ''
   image_name = id+'.tw'
   save_image(image_name)
   with open('./data/'+image_name+'.png','rb') as img_f:
      img_stream = img_f.read()
      img_stream = base64.b64encode(img_stream)
      img_stream = str(img_stream,'utf8')
   return img_stream