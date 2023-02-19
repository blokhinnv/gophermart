import requests
from front.const import LOGIN


class GophermartAPI:
    def __init__(self, base_url):
        self.base_url = base_url

    def logreg(self, login: str, password: str, type: bool) -> dict:
        url = f"{self.base_url}/api/user/{'login' if type == LOGIN else 'register'}"
        r = requests.post(
            url,
            json={"login": login, "password": password},
            headers={"Content-Type": "application/json"},
        )
        if r.status_code == 200:
            return {
                "authorized": True,
                "jwt": r.headers["Authorization"],
                "msg": r.text,
            }
        return {"authorized": False, "jwt": None, "msg": r.text}

    def post_order(self, order_id: str, sess: requests.Session):
        url = f"{self.base_url}/api/user/orders"
        r = sess.post(
            url,
            data=order_id,
            headers={"Content-Type": "text/plain"},
        )
        return {"success": r.status_code == 202, "msg": r.text}

    def withdraw(self, order_id: str, withdraw_sum: str, sess: requests.Session):
        url = f"{self.base_url}/api/user/balance/withdraw"
        r = sess.post(
            url,
            json={"order": order_id, "sum": withdraw_sum},
            headers={"Content-Type": "application/json"},
        )
        return {"success": r.status_code == 200, "msg": r.text}

    def balance(self, sess: requests.Session):
        url = f"{self.base_url}/api/user/balance"
        r = sess.get(url)
        if r.status_code == 200:
            return {"success": True, "response": r.json(), "msg": r.text}
        else:
            return {"success": False, "msg": r.text}

    def get_orders(self, sess: requests.Session):
        url = f"{self.base_url}/api/user/orders"
        r = sess.get(url)

        match r.status_code:
            case 200:
                return {"success": True, "response": r.json(), "msg": r.text}
            case 204:
                return {"success": True, "response": [], "msg": r.text}
            case _:
                return {"success": False, "msg": r.text}
