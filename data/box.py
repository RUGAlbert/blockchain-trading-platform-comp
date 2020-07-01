import numpy as np 
from matplotlib import pyplot as plt 
from numpy import genfromtxt
import seaborn as sns
import pandas as pd

pd.set_option('display.max_columns', 500)

def getSubDf(df, val):
	dfsub = df[['size','sort',val]]
	dfsub = dfsub.rename(columns={val: "values", 'size':"size", "sort":"sort"})
	return dfsub

def getsubs(df):
	dfs = []
	for i in range(0,5):
		dfs.append(getSubDf(df, i))
		
	return dfs

def getDF(name, c_name):
	
	df=pd.read_csv(name, sep=',',header=None)
	df.values
	#df['avg'] = df.mean(axis=1)
	#df['min'] = df.min(axis=1)
	#df['max'] = df.max(axis=1)
	df.index =  [2,10,100,1000]
	df = df.T
	df['parameter'] = c_name
	return df

def createPlot(name, sizes):
	dfs = []
	for i in sizes:
		filename = name+str(i)+'.csv'
		dfs.append(getDF(filename, str(i)))
	index = [2, 10, 100, 1000]
	comb = pd.concat(dfs)
	print(comb)
	df_long = pd.melt(comb, "parameter", var_name="size", value_name="values")
	ax = sns.boxplot(x="size", hue="parameter", y="values", data=df_long)
	ax.set_xlabel("Amount of offers/bids")
	ax.set_ylabel("Time in seconds")
	plt.show()

def createComp(eth, hl):
	df_eth = getDF(eth, 'ethereum')
	df_hl = getDF(hl, 'hyperledger')
	index = [2, 10, 100, 1000]
	columns = ['ethereum', 'hyperledger']
	comb = pd.concat([df_eth, df_hl])
	print(comb)
	df_long = pd.melt(comb, "parameter", var_name="size", value_name="values")
	ax = sns.boxplot(x="size", hue="parameter", y="values", data=df_long)
	ax.set_xlabel("Amount of offers/bids")
	ax.set_ylabel("Time in seconds")
	plt.show()

createPlot('eth_par_',[10, 100, 1000, 10000])
createPlot('hl_par_',[10, 100, 1000])
createPlot('hl_seq_',[10, 100, 1000])
createComp('eth_par_100.csv','hl_par_100.csv')
