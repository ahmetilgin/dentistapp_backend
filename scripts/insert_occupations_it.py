from deep_translator import GoogleTranslator
from pymongo import MongoClient

client = MongoClient("mongodb+srv://admin:wqDohhdGEHR2lI6O@cluster0.s8igoep.mongodb.net/")
db = client["testdb"]
collection = db["professions"]

translator = GoogleTranslator(source='en', target='it')

professions_list = []

# txt dosyasını satır satır oku
with open("occupations.txt", "r", encoding="utf-8") as file:
    for line in file:
        occupation = line.strip()  # Her satırdaki mesleği al
        if occupation:  # Eğer satır boş değilse
            translated = translator.translate(occupation)  # Çeviri yap
            
            profession = {
                "name": translated,
                "count": 0
            }
            professions_list.append(profession)

# Veritabanına eklemek için veri yapısı oluştur
professions_data = {
    "code": "IT",  # Ülke kodunu belirt
    "professions": professions_list
}

# Veriyi MongoDB'ye ekle
collection.insert_one(professions_data)

print("Veriler başarıyla çevrildi ve MongoDB'ye uygun yapıda eklendi.")
