import os
import front.pages.auth
import front.pages.home
import front.pages.not_found
import front.pages.register
from dash import dcc, html
from dash_extensions.enrich import Input, Output, html
from dotenv import load_dotenv
from flask_login import current_user
from front.api.gophermart import GophermartAPI
from front.app import app
from loguru import logger

app.layout = html.Div(
    [
        html.Div(id="page-content", className="content"),
        html.Div(
            [
                dcc.Location(id="url", refresh=False),
            ],
            id="url_container",
        ),
    ]
)


@app.callback(
    Output("page-content", "children"),
    Output("url_container", "children"),
    Input("url", "pathname"),
    prevent_initial_call=True,
)
def display_page(pathname):
    def new_loc(url):
        return dcc.Location(id="url", pathname=url, refresh=False)

    to_home = (front.pages.home.layout, new_loc("/"))
    to_login = (front.pages.auth.layout, new_loc("/login"))
    to_register = (front.pages.register.layout, new_loc("/register"))
    to_not_found = (front.pages.not_found.layout, new_loc("/404"))

    if pathname == "/":
        if current_user.is_authenticated:
            return to_home
        return to_login
    elif pathname == "/login":
        if current_user.is_authenticated:
            return to_home
        return to_login
    elif pathname == "/register":
        if current_user.is_authenticated:
            return to_home
        return to_register
    else:
        return to_not_found


if __name__ == "__main__":
    load_dotenv()
    logger.info("Started")
    api = GophermartAPI(os.getenv("GOPHERMART_API_ADDR"))
    app.api = api
    app.USERS = {}
    app.run_server(debug=True, port="7777")
