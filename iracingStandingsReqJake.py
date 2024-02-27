import json
import requests
import math

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

param_sets = []
with open ('2-23-2024-jake-standings-input.json', 'r') as input_file:
    json_data = json.load(input_file)
    for result in json_data:
        param_sets.append(result.split('_'))

query_string = 'https://members-ng.iracing.com/data/stats/season_driver_standings?season_id={}&car_class_id={}'
division_query_string = 'https://members-ng.iracing.com/data/stats/season_driver_standings?season_id={}&car_class_id={}&division={}'
output = []
for index, param_set in enumerate(param_sets):
  print(str(index + 1) + ' of ' + str(param_sets.__len__()))
  season_dict = {}
  [season_id, car_class_id] = param_set
  response = session.get(query_string.format(season_id, car_class_id))
  response_json = response.json()
  standings_response = session.get(response_json['link'])
  standings_response_json = standings_response.json()
  my_rank = standings_response_json['customer_rank']
  calculated_page_number = my_rank / standings_response_json['chunk_info']['chunk_size']
  page_number = calculated_page_number if (calculated_page_number % 1 != 0) else (calculated_page_number - 1)
  root_url = standings_response_json['chunk_info']['base_download_url']
  season_name = standings_response_json['season_name']
  season_dict['season_name'] = season_name
  season_dict['overall_rank'] = my_rank
  season_dict['season_id'] = season_id
  season_dict['car_class_id'] = car_class_id
  if (my_rank > 0):
      file_name = standings_response_json['chunk_info']['chunk_file_names'][math.floor(page_number)]
      file_response = session.get(root_url + file_name)
      file_response_json = file_response.json()
      [me] = [driver for driver in file_response_json if driver['cust_id'] == 300752]
      my_division = me['division']
      season_dict['division'] = my_division + 1
      season_dict['season_driver_data'] = me
      if (my_division > -1):
        division_response = session.get(division_query_string.format(season_id, car_class_id, my_division))
        division_response_json = division_response.json()
        division_standings_response = session.get(division_response_json['link'])
        division_standings_response_json = division_standings_response.json()
        my_division_rank = division_standings_response_json['customer_rank']
        season_dict['division_rank'] = my_division_rank
      output.append(season_dict)
with open('output.json', 'w') as output_file:
  json.dump(output, output_file)
