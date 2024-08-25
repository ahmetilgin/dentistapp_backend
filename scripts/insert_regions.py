import requests
import json
from bson import ObjectId
from pymongo import MongoClient



# Convert to dictionary for JSON serialization
def obj_to_dict(obj):
    if isinstance(obj, ObjectId):
        return str(obj)
    elif hasattr(obj, '__dict__'):
        return {key: obj_to_dict(value) for key, value in obj.__dict__.items()}
    elif isinstance(obj, list):
        return [obj_to_dict(item) for item in obj]
    else:
        return obj

def get_country_data(country_name, country_code):
    # Overpass API endpoint
    overpass_url = "http://overpass-api.de/api/interpreter"

    # Overpass QL query for cities (provinces)
    city_query = f"""
    [out:json];
    area["name"="{country_name}"]->.country;
    relation(area.country)["admin_level"="4"]["boundary"="administrative"];
    out center;
    """

    # Overpass QL query for districts
    district_query = """
    [out:json];
    area["name"="{}"]->.city;
    relation(area.city)["admin_level"="6"]["boundary"="administrative"];
    out center;
    """

    # Send request to Overpass API for cities
    response = requests.get(overpass_url, params={'data': city_query})
    city_data = response.json()
    filename = f"{country_name.lower().replace(' ', '_')}_raw.json"
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(city_data, f, ensure_ascii=False, indent=2, default=obj_to_dict)

    class District:
        def __init__(self, name, latitude, longitude):
            self.name = name
            self.latitude = latitude
            self.longitude = longitude

    class City:
        def __init__(self, name, districts, latitude, longitude):
            self.name = name
            self.districts = districts
            self.latitude = latitude
            self.longitude = longitude

    class Country:
        def __init__(self,  name, code, cities):
            self.name = name
            self.code = code
            self.cities = cities
        def to_dict(self):
            return obj_to_dict(self)

    # Process the city data
    country = {
        "name": country_name,
        "code": country_code,
        "cities": []
    }

    for element in city_data['elements']:
        if element['type'] == 'relation':
            tags = element.get('tags', {})
            name = tags.get('name', '')
            center = element.get('center', {})
            latitude = center.get('lat', 0)
            longitude = center.get('lon', 0)
            print(name)
            # Sadece ismi olan şehirleri ekleyin
            if name:
                # Get districts for this city
                district_response = requests.get(overpass_url, params={'data': district_query.format(name)})
                district_data = district_response.json()

                districts = []
                for district_element in district_data['elements']:
                    if district_element['type'] == 'relation':
                        district_tags = district_element.get('tags', {})
                        district_name = district_tags.get('name', '')
                        district_center = district_element.get('center', {})
                        district_latitude = district_center.get('lat', 0)
                        district_longitude = district_center.get('lon', 0)
                        if district_name:
                            districts.append(District(district_name, district_latitude, district_longitude))

                country['cities'].append(City(name, districts, latitude, longitude))

    # Create Country object
    country_obj = Country(
        name=country['name'],
        code=country['code'],
        cities=country['cities']
    )

    # Sort cities and districts alphabetically
    country_obj.cities.sort(key=lambda x: x.name)
    for city in country_obj.cities:
        city.districts.sort(key=lambda x: x.name)



    # Save the data to a JSON file
    filename = f"{country_name.lower().replace(' ', '_')}_cities_districts.json"
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(country_obj, f, ensure_ascii=False, indent=2, default=obj_to_dict)

    print(f"Country: {country_obj.name} ({country_obj.code})")
    print(f"Country: {country_obj.name} ({country_obj.code})")
    print(f"Total cities: {len(country_obj.cities)}")
    print(f"Total districts: {sum(len(city.districts) for city in country_obj.cities)}")
    print(f"Data has been processed and saved to '{filename}'.")

    return country_obj


def get_albania_data():
    overpass_url = "http://overpass-api.de/api/interpreter"
    
    # County (qark) sorgusu
    county_query = """
    [out:json];
    area["name"="Shqipëria"]->.albania;
    relation(area.albania)["admin_level"="6"]["boundary"="administrative"];
    out center;
    """
    
    # Municipality (bashki) sorgusu
    municipality_query = """
    [out:json];
    area["name"="{}"]->.county;
    relation(area.county)["admin_level"="8"]["boundary"="administrative"];
    out center;
    """
    
    response = requests.get(overpass_url, params={'data': county_query})
    county_data = response.json()
    
    albania = {
        "name": "Albania",
        "code": "AL",
        "counties": []
    }
    
    for element in county_data['elements']:
        if element['type'] == 'relation':
            tags = element['tags']
            name = tags.get('name', '')
            center = element.get('center', {})
            latitude = center.get('lat', 0)
            longitude = center.get('lon', 0)
            
            # Get districts for this county
            municipality_response = requests.get(overpass_url, params={'data': municipality_query.format(name)})
            municipality_data = municipality_response.json()
            
            districts = []
            for mun_element in municipality_data['elements']:
                if mun_element['type'] == 'relation':
                    mun_tags = mun_element['tags']
                    mun_name = mun_tags.get('name', '')
                    mun_center = mun_element.get('center', {})
                    mun_latitude = mun_center.get('lat', 0)
                    mun_longitude = mun_center.get('lon', 0)
                    districts.append({"name": mun_name, "latitude": mun_latitude, "longitude": mun_longitude})
            
            albania['counties'].append({
                "name": name,
                "latitude": latitude,
                "longitude": longitude,
                "districts": districts
            })
    country_name = "albania"
    filename = f"{country_name.lower().replace(' ', '_')}_cities_districts.json"
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(albania, f, ensure_ascii=False, indent=2, default=obj_to_dict)

    return albania


def send_to_mongo(country_data):
    client = MongoClient('mongodb+srv://admin:wqDohhdGEHR2lI6O@cluster0.s8igoep.mongodb.net/')
    db = client['testdb']
    collection = db['country']
    country_dict = obj_to_dict(country_data)
    collection.insert_one(country_dict)

if __name__ == "__main__":
    # send_to_mongo(get_country_data("Türkiye", "TR"))
    # send_to_mongo(get_country_data("Shqipëria", "ALB"))
    send_to_mongo(get_albania_data())


