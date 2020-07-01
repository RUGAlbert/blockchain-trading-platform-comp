import numpy as np 
from matplotlib import pyplot as plt 
from numpy import genfromtxt
import pandas as pd

pd.set_option('display.max_columns', 500)

def getDF(name):
	
	df=pd.read_csv(name, sep=',',header=None)
	df.values
	df['avg'] = df.mean(axis=1)
	df['min'] = df.min(axis=1)
	df['max'] = df.max(axis=1)
	df['size'] = [2, 10, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]
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
	#columns = []
	comb = pd.DataFrame(index=index, columns=columns)
	for i in range(len(sizes)):
		comb[str(sizes[i])] = dfs[i]['avg']
		#comb[str(sizes[i]) + "_avg"] = dfs[i]['avg']
		#comb[str(sizes[i]) + "_min"] = dfs[i]['min']
		#comb[str(sizes[i]) + "_max"] = dfs[i]['max']
	print(comb)
	ax = comb.plot(y=columns, kind='line')
	ax.set_xlabel("Amount of offers/bids")
	ax.set_ylabel("Time in seconds")
	plt.show()

def createComp(eth, hl):
	df_eth = getDF(eth)
	df_hl = getDF(hl)
	index = [2, 10, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]
	columns = ['ethereum', 'hyperledger']
	#columns = []
	comb = pd.DataFrame(index=index, columns=columns)
	comb['ethereum'] = df_eth['avg']
	comb['hyperledger'] = df_hl['avg']
	
	#comb['eth_min'] = df_eth['min']
	#comb['eth_avg'] = df_eth['avg']
	#comb['eth_max'] = df_eth['max']
	#comb['hl_min'] = df_hl['min']
	#comb['hl_avg'] = df_hl['avg']
	#comb['hl_max'] = df_hl['max']
	#print(comb)
	ax = comb.plot(y=columns, kind='line')
	ax.set_xlabel("Amount of offers/bids")
	ax.set_ylabel("Time in seconds")
	plt.show()

#createPlot('eth_par_',[10, 100, 1000, 10000])
#createPlot('hl_par_',[10, 100, 1000])
#createPlot('hl_seq_',[10, 100, 1000])
#createComp('eth_par_100.csv','hl_par_100.csv')
createComp('eth_scalability.csv','hl_scalability.csv')
