import numpy as np 
from matplotlib import pyplot as plt 
from numpy import genfromtxt
import pandas as pd

def getDF(name):
	
	df=pd.read_csv(name, sep=',',header=None)
	df.values
	df['avg'] = df.mean(axis=1)
	df['min'] = df.min(axis=1)
	df['max'] = df.max(axis=1)
	df['size'] = [2,10,100,1000]
	df.index = df['size']
	return df

def createPlot(name, sizes):
	dfs = []
	for i in sizes:
		filename = name+str(i)+'.csv'
		print(filename)
		dfs.append(getDF(filename))
	index = [2, 10, 100, 1000]
	columns = list(map(str, sizes))
	comb = pd.DataFrame(index=index, columns=columns)
	for i in range(len(sizes)):
		comb[str(sizes[i])] = dfs[i]['avg']
	print(comb)
	comb.plot(y=columns, kind='line')
	plt.show()

def createComp(eth, hl):
	df_eth = getDF(eth)
	df_hl = getDF(hl)
	index = [2, 10, 100, 1000]
	columns = ['eth', 'hl']
	comb = pd.DataFrame(index=index, columns=columns)
	comb['eth'] = df_eth['avg']
	comb['hl'] = df_hl['avg']
	print(comb)
	comb.plot(y=columns, kind='line')
	plt.show()

#createPlot('eth_par_',[10, 100, 1000, 10000])
#createPlot('hl_par_',[10, 100, 1000])
#createPlot('hl_seq_',[10, 100, 1000])
createComp('eth_par_100.csv','hl_par_100.csv')
