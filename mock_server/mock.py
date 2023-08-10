from flask import Flask, request, jsonify
from functools import wraps

app = Flask(__name__)

# Mock data
mock_users = {
    "1": {
        "id": 1,
        "employeeId": 1,
        "firstName": "John",
        "lastName": "Doe",
        "email": "john.doe@bamboohr.com",
        "status": "enabled",
        "lastLogin": "2011-03-19T10:16:00+00:00"
    },
    "2": {
        "id": 2,
        "firstName": "Jane",
        "lastName": "Doe",
        "email": "jane.doe@bamboohr.com",
        "status": "enabled",
        "lastLogin": "2011-08-29T11:17:43+00:00"
    },
    "3": {
        "id": 3,
        "employeeId": 2,
        "firstName": "Michael",
        "lastName": "Smith",
        "email": "michael.smith@bamboohr.com",
        "status": "enabled",
        "lastLogin": "2023-08-01T08:00:00+00:00"
    }
}

# Basic Authentication decorator
def requires_auth(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        auth = request.authorization
        if not auth or auth.username != 'APIKEY' or auth.password != 'x':
            return jsonify(message='Unauthorized'), 401
        return f(*args, **kwargs)
    return decorated

@app.route('/testcompany/v1/meta/users/', methods=['GET'])
@requires_auth
def get_users():
    filtered_users = {user_id: user for user_id, user in mock_users.items() if 'employeeId' in user}
    return jsonify(filtered_users)

if __name__ == '__main__':
    app.run(host='localhost', port=8000)

