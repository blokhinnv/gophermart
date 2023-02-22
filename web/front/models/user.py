import requests
from flask_login import UserMixin


class User(UserMixin):
    def __init__(self, username: str) -> None:
        self.username = username
        self.id = username
        self.sess = None

        super().__init__()

    def start_session(self, jwt: str) -> None:
        self.sess = requests.Session()
        self.sess.headers.update({"Authorization": jwt})
