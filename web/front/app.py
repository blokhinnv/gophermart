import dash_bootstrap_components as dbc
from dash_extensions.enrich import (
    DashProxy,
    MultiplexerTransform,
    NoOutputTransform,
    TriggerTransform,
)
from flask_login import LoginManager

app = DashProxy(
    __name__,
    external_stylesheets=[dbc.themes.BOOTSTRAP],
    transforms=[TriggerTransform(), NoOutputTransform(), MultiplexerTransform()],
    suppress_callback_exceptions=True,
)
server = app.server
login_manager = LoginManager()
login_manager.init_app(server)
login_manager.login_view = "/login"
# callback to reload the user object
@login_manager.user_loader
def load_user(user_id):
    return app.USERS.get(user_id)
