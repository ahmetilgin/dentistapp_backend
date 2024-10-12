from pymongo import MongoClient

client = MongoClient("mongodb+srv://admin:wqDohhdGEHR2lI6O@cluster0.s8igoep.mongodb.net/")
db = client["testdb"]
collection = db["professions"]

professions_list = []

with open("occupations.txt", "r", encoding="utf-8") as file:
    for line in file:
        occupation = line.strip()  # Her satırdaki mesleği al
        if occupation:  # Eğer satır boş değilse
            profession = {
                "name": occupation,
                "count": 0
            }
            professions_list.append(profession)

professions_data = {
    "code": "EN",  # Buraya uygun ülke kodunu ekle
    "professions": professions_list
}

# Veriyi MongoDB'ye ekle
collection.insert_one(professions_data)

print("Veriler başarıyla çevrildi ve MongoDB'ye uygun yapıda eklendi.")
