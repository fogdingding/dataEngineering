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

#這邊為boyer_moore的演算法 字串比對。
def string_match_boyer_moore(string,match,start=0):
    string_len = len(string)
    match_len  = len(match)
    end = match_len - 1
    if string_len < match_len:
        return start
    while string[end] == match[end]:
        end -= 1
        if end == 0:
            return ('yes')
    idx = contain_char(match,string[end])
    shift = match_len
    if idx > -1:
        shift = end - idx
    start += shift
    string_match_boyer_moore(string[shift:],match,start)

#這邊負責計算字元相等否。
def contain_char(s,c):
   for i in range(len(s)):
      if c == s[i]:
          return i
   return -1

#這邊負責呼叫grep函數
def grep(file_path,match_string):
    app=[]
    with open(file_path,encoding="utf-8") as f_ssv:
        for line in f_ssv:
            line_string = line
            if (string_match_boyer_moore(line_string,match_string) == 'yes'):
                app.append(line.strip())
    return app

@app.route('/')
# 這邊為主畫面
@app.route('/index')
def index():
   # render_template 內建可引導至templates資料夾的index.html，並傳送title='資料工程'。
   return render_template('index.html' , title='資料工程')

#這邊為搜尋完的頁面
@app.route('/search/<search_text>')
def search(search_text):
   #這邊為自製的grep
   output_list = grep('.\data\output.ssv',search_text)
   return render_template('search.html',search_text = output_list)

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