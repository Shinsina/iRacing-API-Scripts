import json
import requests

customer_id_to_season_quarter_year_map = {
  "300752": [
    { "season_quarter": 4, "season_year": 2017 },
    { "season_quarter": 1, "season_year": 2018 },
    { "season_quarter": 4, "season_year": 2018 },
    { "season_quarter": 1, "season_year": 2019 },
    { "season_quarter": 2, "season_year": 2019 },
    { "season_quarter": 3, "season_year": 2019 },
    { "season_quarter": 1, "season_year": 2020 },
    { "season_quarter": 2, "season_year": 2020 },
    { "season_quarter": 3, "season_year": 2020 },
    { "season_quarter": 1, "season_year": 2022 },
    { "season_quarter": 2, "season_year": 2022 },
    { "season_quarter": 3, "season_year": 2022 },
    { "season_quarter": 4, "season_year": 2022 },
    { "season_quarter": 1, "season_year": 2023 },
    { "season_quarter": 2, "season_year": 2023 },
    { "season_quarter": 3, "season_year": 2023 },
    { "season_quarter": 4, "season_year": 2023 },
    { "season_quarter": 1, "season_year": 2024 },
    { "season_quarter": 2, "season_year": 2024 },
  ],
  "815162": [
    { "season_quarter": 3, "season_year": 2022 },
    { "season_quarter": 4, "season_year": 2022 },
    { "season_quarter": 1, "season_year": 2023 },
    { "season_quarter": 2, "season_year": 2023 },
    { "season_quarter": 3, "season_year": 2023 },
    { "season_quarter": 4, "season_year": 2023 },
    { "season_quarter": 1, "season_year": 2024 },
    { "season_quarter": 2, "season_year": 2024 },
  ],
};

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

output = {}
query_string = 'https://members-ng.iracing.com/data/results/search_series?season_quarter={}&season_year={}&cust_id={}&official_only=true&event_types=5'
cust_ids = list(customer_id_to_season_quarter_year_map.keys())
for index, cust_id in enumerate(cust_ids):
   param_sets = customer_id_to_season_quarter_year_map[cust_id]
   for index, param_set in enumerate(param_sets):
      response = session.get(query_string.format(param_set['season_quarter'], param_set['season_year'], cust_id))
      response_json = response.json()
      print(query_string.format(param_set['season_quarter'], param_set['season_year'], cust_id))
      base_download_url = response_json['data']['chunk_info']['base_download_url']
      chunk_file_names = response_json['data']['chunk_info']['chunk_file_names']
      for chunk_file_name in chunk_file_names:
        chunk_file_response = requests.get(base_download_url + chunk_file_name)
        chunk_file_response_json = chunk_file_response.json()
        output[chunk_file_name] = chunk_file_response_json

with open('search-series-output.json', 'w') as output_file:
  json.dump(output, output_file)
