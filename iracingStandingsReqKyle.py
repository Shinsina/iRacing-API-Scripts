import json
import requests

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

param_sets = []
with open ('7-18-2024-kyle-standings-input.json', 'r') as input_file:
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
  root_url = standings_response_json['chunk_info']['base_download_url']
  season_name = standings_response_json['season_name']
  season_dict['season_name'] = season_name
  season_dict['season_id'] = season_id
  season_dict['car_class_id'] = car_class_id
  for file_name in standings_response_json['chunk_info']['chunk_file_names']:
    file_response = session.get(root_url + file_name)
    # @todo figure out why this 403's
    if (file_response.status_code == 200):
      file_response_json = file_response.json()
      potential_kyles = [driver for driver in file_response_json if driver['cust_id'] == 251134]
      if (len(potential_kyles) > 0):
        [kyle] = potential_kyles
        kyle_division = kyle['division']
        kyle_rank = kyle['rank']
        season_dict['division'] = kyle_division + 1
        season_dict['season_driver_data'] = kyle
        season_dict['overall_rank'] = kyle_rank
        division_response = session.get(division_query_string.format(season_id, car_class_id, kyle_division))
        division_response_json = division_response.json()
        division_standings_response = session.get(division_response_json['link'])
        division_standings_response_json = division_standings_response.json()
        root_url = division_standings_response_json['chunk_info']['base_download_url']
        for file_name in division_standings_response_json['chunk_info']['chunk_file_names']:
            file_response = session.get(root_url + file_name)
            file_response_json = file_response.json()
            potential_kyles = [driver for driver in file_response_json if driver['cust_id'] == 251134]
            if (len(potential_kyles) > 0):
              [kyle] = potential_kyles
              kyle_division_rank = kyle['rank']
              season_dict['division_rank'] = kyle_division_rank
  output.append(season_dict)
with open('kyle-standings-output.json', 'w') as output_file:
  json.dump(output, output_file)
