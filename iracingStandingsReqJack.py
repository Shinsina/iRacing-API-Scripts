import json
import requests

# @todo Repull this list
# Array.from(new Set(subsessions.map((subsession) => subsession.session_id))).forEach((item) => console.log(item + ','));
# In pages/user/[id]/subsessions.astro
season_id_list = [
   4150,
4151,
4140,
4137,
4181,
4155,
4139,
4077,
3998,
3999,
3996,
3980,
3984,
4014,
4092,
3862,
3861,
3875,
3850,
3848,
3847,
3846,
3946,
3947,
3725,
3790,
4150,
4151,
4140,
4137,
4181,
4155,
4139,
4077,
3998,
3999,
3996,
3980,
3984,
4014,
4092,
3862,
3861,
3875,
3850,
3848,
3847,
3846,
3946,
3947,
3725,
3790,
4150,
4151,
4140,
4137,
4181,
4155,
4139,
4077,
3998,
3999,
3996,
3980,
3984,
4014,
4092,
3862,
3861,
3875,
3850,
3848,
3847,
3846,
3946,
3947,
3725,
3790,
]

session = requests.session()

with open('cookie.txt', 'r') as file:
    cookies = requests.utils.cookiejar_from_dict(json.load(file))
    session.cookies.update(cookies)

param_sets = []
with open ('input.json', 'r') as input_file:
    json_data = json.load(input_file)
    for result in json_data:
        param_sets.append(result.split('_'))

query_string = 'https://members-ng.iracing.com/data/stats/season_driver_standings?season_id={}&car_class_id={}'
division_query_string = 'https://members-ng.iracing.com/data/stats/season_driver_standings?season_id={}&car_class_id={}&division={}'
output = []
for param_set in param_sets:
  season_dict = {}
  [season_id, car_class_id] = param_set
  if int(season_id) in season_id_list:
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
        potential_jacks = [driver for driver in file_response_json if driver['cust_id'] == 815162]
        if (len(potential_jacks) > 0):
          [jack] = potential_jacks
          jack_division = jack['division']
          jack_rank = jack['rank']
          season_dict['division'] = jack_division + 1
          season_dict['season_driver_data'] = jack
          season_dict['overall_rank'] = jack_rank
          division_response = session.get(division_query_string.format(season_id, car_class_id, jack_division))
          division_response_json = division_response.json()
          division_standings_response = session.get(division_response_json['link'])
          division_standings_response_json = division_standings_response.json()
          root_url = division_standings_response_json['chunk_info']['base_download_url']
          for file_name in division_standings_response_json['chunk_info']['chunk_file_names']:
              file_response = session.get(root_url + file_name)
              file_response_json = file_response.json()
              potential_jacks = [driver for driver in file_response_json if driver['cust_id'] == 815162]
              if (len(potential_jacks) > 0):
                [jack] = potential_jacks
                jack_division_rank = jack['rank']
                season_dict['division_rank'] = jack_division_rank
    output.append(season_dict)
with open('output.json', 'w') as output_file:
  json.dump(output, output_file)
