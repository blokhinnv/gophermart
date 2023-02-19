from dash import html
from front.components.logreg_aio import LogRegAIO
from front.const import LOGIN

layout = html.Div(LogRegAIO(type=LOGIN))
