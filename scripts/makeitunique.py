# Girdi dosyasının adı
input_file_name = 'occupations.txt'  # Orijinal listeyi içeren dosyanın adı
# Çıktı dosyasının adı
output_file_name = 'unique_professions.txt'  # Benzersiz satırların yazılacağı dosyanın adı

def get_unique_lines(input_file, output_file):
    unique_lines = set()  # Benzersiz satırları saklamak için bir set

    # Girdi dosyasını oku
    with open(input_file, 'r') as infile:
        for line in infile:
            unique_lines.add(line.strip())  # Satırı set'e ekle (strip ile boşlukları temizle)

    # Benzersiz satırları çıktı dosyasına yaz
    with open(output_file, 'w') as outfile:
        for line in unique_lines:
            outfile.write(line + '\n')  # Her satırı yeni bir satıra yaz

# Fonksiyonu çağır
get_unique_lines(input_file_name, output_file_name)
print(f"Benzersiz satırlar '{output_file_name}' dosyasına yazıldı.")