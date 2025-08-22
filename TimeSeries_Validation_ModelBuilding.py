# Importing the basic preprocessing packages
import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
# The series module from Pandas will help in creating a time series
from pandas import Series,DataFrame
import seaborn as sns
#%matplotlib inline

# About the Data Set (Location: https://www.kaggle.com/sumanthvrao/daily-climate-time-series-data) 
# To forecast the daily climate of a city in India
time_series = pd.read_csv('data/DailyDelhiClimateTrain.csv', parse_dates=['date'], index_col='date')



time_series.head()


time_series.plot(subplots=True)

plt.show()