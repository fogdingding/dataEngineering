#!/usr/bin/env python
# coding: utf-8

# In[16]:


def string_match_boyer_moore(string,match,start=0):
    string_len = len(string)
    match_len  = len(match)
    end = match_len - 1
    if string_len < match_len:
        return start;
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

def contain_char(s,c):
   for i in range(len(s)):
      if c == s[i]:
          return i
   return -1


# In[34]:


def grep(file_path,match_string):
    app=[]
    with open(file_path,encoding="utf-8") as f_ssv:
        for line in f_ssv:
            line_string = line
            if (string_match_boyer_moore(line_string,match_string) == 'yes'):
                app.append(line.strip())
    return app


# In[35]:


a = grep('./output.ssv','110')


# In[36]:


a

