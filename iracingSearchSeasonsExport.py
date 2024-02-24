import json
import requests

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

past_series = []
with open ('past-season-series-ids-input.json', 'r') as input_file:
    json_data = json.load(input_file)
    past_series.extend(json_data)

output = {}
query_string = 'https://members-ng.iracing.com/data/results/search_series?season_quarter=1&season_year=2024&cust_id=300752&official_only=true&event_types=5'
response = session.get(query_string)
response_json = response.json()
base_download_url = response_json['data']['chunk_info']['base_download_url']
chunk_file_names = response_json['data']['chunk_info']['chunk_file_names']
for chunk_file_name in chunk_file_names:
   chunk_file_response = requests.get(base_download_url + chunk_file_name)
   chunk_file_response_json = chunk_file_response.json()
   output[chunk_file_name] = chunk_file_response_json

with open('output.json', 'w') as output_file:
  json.dump(output, output_file)
