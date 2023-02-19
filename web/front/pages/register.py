from dash import html
from front.components.logreg_aio import LogRegAIO
from front.const import REGISTER

layout = html.Div(LogRegAIO(type=REGISTER))
