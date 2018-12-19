#!/usr/bin/env python
# -*- coding: utf-8 -*
from bs4 import BeautifulSoup
import requests

result = requests.get("http://www.tej.com.tw/webtej/doc/uid.htm#上市")
result.encoding = 'big5'
soup = BeautifulSoup(result.text, 'html.parser')

app=[]
for link in soup.find_all('table'):
    app.append(link)
#1、24、47、50

app2=[]
for link in app[1].find_all('tr'):
    for link2 in link.find_all('td'):
        if( link2.get_text() ):
            if ( len(link2.get_text().split()) < 3 and len(link2.get_text().split()) > 1):
                app2.append(link2.get_text().split())
for link in app[24].find_all('tr'):
    for link2 in link.find_all('td'):
        if( link2.get_text() ):
            if ( len(link2.get_text().split()) < 3 and len(link2.get_text().split()) > 1):
                app2.append(link2.get_text().split())
for link in app[47].find_all('tr'):
    for link2 in link.find_all('td'):
        if( link2.get_text() ):
            if ( len(link2.get_text().split()) < 3 and len(link2.get_text().split()) > 1):
                app2.append(link2.get_text().split())
for link in app[50].find_all('tr'):
    for link2 in link.find_all('td'):
        if( link2.get_text() ):
            if ( len(link2.get_text().split()) < 3 and len(link2.get_text().split()) > 1):
                app2.append(link2.get_text().split())

app3=[]
for i in app2:
    if not i in app3:
        app3.append(i)

import os
os.system('ls')
os.system('grep 1 test.txt -n -w > test.tsv')

import csv

with open('output.tsv', 'w', newline='' , encoding="utf-8") as f:
    writer = csv.writer(f,delimiter='\t')
    writer.writerows(app3)