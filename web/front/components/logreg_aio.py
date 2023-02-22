import uuid

import dash_bootstrap_components as dbc
import dash_loading_spinners as dls
from dash import dcc, html, no_update
from dash_extensions.enrich import ALL, MATCH, Input, Output, State, Trigger, callback
from flask_login import login_user
from front.const import LOGIN
from front.models.user import User
from front.app import app


# All-in-One Components should be suffixed with 'AIO'
class LogRegAIO(html.Div):  # html.Div will be the "parent" component
    # A set of functions that create pattern-matching callbacks of the subcomponents
    class ids:
        login = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "login_input",
            "aio_id": aio_id,
        }
        password = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "password_input",
            "aio_id": aio_id,
        }
        alert = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "alerts",
            "aio_id": aio_id,
        }
        button = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "button",
            "aio_id": aio_id,
        }
        location = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "location",
            "aio_id": aio_id,
        }
        type = lambda aio_id: {
            "component": "LogRegAIO",
            "subcomponent": "type",
            "aio_id": aio_id,
        }

    # Make the ids class a public class
    ids = ids

    # Define the arguments of the All-in-One component
    def __init__(self, type: str, aio_id=None):
        if aio_id is None:
            aio_id = str(uuid.uuid4())

        super().__init__(
            dls.Hash(
                html.Div(
                    [
                        html.Div(
                            [
                                dcc.Store(id=self.ids.type(aio_id), data=type),
                                dcc.Store(id=self.ids.location(aio_id), data=None),
                                html.H1(
                                    children="Форма {}".format(
                                        "авторизации"
                                        if type == LOGIN
                                        else "регистрации"
                                    )
                                ),
                                html.Div(
                                    id=self.ids.alert(aio_id),
                                ),
                                dbc.Form(
                                    [
                                        dbc.Row(
                                            [
                                                dbc.Label(
                                                    "Логин",
                                                    html_for=self.ids.login(aio_id),
                                                    width=2,
                                                ),
                                                dbc.Col(
                                                    dbc.Input(
                                                        type="input",
                                                        id=self.ids.login(aio_id),
                                                        placeholder="Введите логин",
                                                    ),
                                                    width=10,
                                                ),
                                            ],
                                            className="mb-3",
                                        ),
                                        dbc.Row(
                                            [
                                                dbc.Label(
                                                    "Пароль",
                                                    html_for=self.ids.password(aio_id),
                                                    width=2,
                                                ),
                                                dbc.Col(
                                                    dbc.Input(
                                                        type="password",
                                                        id=self.ids.password(aio_id),
                                                        placeholder="Введите пароль",
                                                    ),
                                                    width=10,
                                                ),
                                            ],
                                            className="mb-3",
                                        ),
                                    ]
                                ),
                                dbc.Button(
                                    "Авторизация" if type == LOGIN else "Регистрация",
                                    color="primary",
                                    className="me-1",
                                    id=self.ids.button(aio_id),
                                ),
                            ],
                            className="w-30",
                        )
                    ],
                    className="d-flex justify-content-center align-items-center",
                    style={"height": "100vh"},
                ),
                color="#435278",
                speed_multiplier=2,
                size=100,
                debounce=200,
                fullscreen=True,
            )
        )

    @callback(
        Output(ids.location(MATCH), "data"),
        Output(ids.alert(MATCH), "children"),
        Trigger(ids.button(MATCH), "n_clicks"),
        State(ids.login(MATCH), "value"),
        State(ids.password(MATCH), "value"),
        State(ids.type(MATCH), "data"),
        prevent_initial_call=True,
    )
    def process(login, password, type):
        if login is None or password is None:
            return no_update, dbc.Alert(
                "Оба поля обязательны к заполнению",
                color="danger",
                duration=3000,
            )

        res = app.api.logreg(login, password, type)
        if res["authorized"]:
            user = User(login)
            user.start_session(res["jwt"])
            app.USERS[user.id] = user
            login_user(user)
            return "/", no_update

        return no_update, dbc.Alert(
            res["msg"],
            color="danger",
            duration=3000,
        )


@callback(
    Output("url", "pathname"),
    Input(LogRegAIO.ids.location(ALL), "data"),
    prevent_initial_call=True,
)
def redirect(urls):
    if not urls:
        return no_update
    url = urls[0]
    if not url:
        return no_update
    return url
