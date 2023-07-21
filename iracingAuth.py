import hashlib
import base64
import requests
import json
import os

# Needs to be executed via auth.sh
# Additionally IRACING_USERNAME and IRACING_PASSWORD must be set in a .env file at root

username = os.environ.get('IRACING_USERNAME')
password = os.environ.get('IRACING_PASSWORD')

def encode_pw(username, password):
    initialHash = hashlib.sha256((password + username.lower()).encode('utf-8')).digest()

    hashInBase64 = base64.b64encode(initialHash).decode('utf-8')

    return hashInBase64


pwValueToSubmit = encode_pw(username, password)


auth_response = requests.post('https://members-ng.iracing.com/auth', data=None,json={"email": username, "password": pwValueToSubmit}, headers={"Content-Type": "application/json"})

session = requests.session()

with open('cookie.txt', 'w') as file:
    json.dump(requests.utils.dict_from_cookiejar(auth_response.cookies), file)
