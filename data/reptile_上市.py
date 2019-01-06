#!/usr/bin/env python
# coding: utf-8

# In[5]:


# -*- coding: utf-8 -*
from bs4 import BeautifulSoup
import requests


# In[6]:


def find_html_td_text(bs4_html_tr):
    bs4_html_td_list = []
    app = bs4_html_tr.find_all('td')
    for i in range (len(app)):
        bs4_html_td_list.append(app[i].get_text().split())
    return bs4_html_td_list


# In[7]:


def get_array_list_to_string(array_list):
    test = array_list
    string=""
    for i in range(len(test)):
        locals()["str%s" %i] = ' '.join(test[i])
    for i in range(len(test)):
        if( i == len(test)-1):
            string += locals()["str%s" %i]
        elif( test[i] == []):
            string += locals()["str%s" %i]
        else:
            string += locals()["str%s" %i]+" "
    return string


# In[8]:


result = requests.get("http://isin.twse.com.tw/isin/C_public.jsp?strMode=2")
result.encoding = 'big5'
soup = BeautifulSoup(result.text, 'html.parser')


# In[9]:


app=[]
for link in soup.find_all('table'):
    for index, link2 in enumerate( link.find_all('tr') ):
        app.append(find_html_td_text(link2))
        if (app[index] == [['上市認購(售)權證']]):
            break
#931
#929


# In[53]:


app2=[]
for link in app:
    app2.append(get_array_list_to_string(link))
app2.pop(-1)
app2.pop(1)


# In[55]:


import csv

with open('output.ssv', 'w', newline='' , encoding="utf-8") as f:
    for item in app2:
        f.write("%s\n" % item)


# In[ ]:




