import pandas as pd
from pymongo import MongoClient

client = MongoClient("mongodb+srv://admin:wqDohhdGEHR2lI6O@cluster0.s8igoep.mongodb.net/")
db = client["testdb"]
collection = db["professions"]

df = pd.read_csv("occupations.csv", header=None, names=["level", "code", "occupation", "active"])

occupations = df[["code", "occupation"]]

for index, row in occupations.iterrows():
    if pd.notna(row["occupation"]):  
        profession = {
            "name": row["occupation"],
            "search_counter": 0
        }
        collection.insert_one(profession)

print("Veriler başarıyla MongoDB'ye eklendi.")
