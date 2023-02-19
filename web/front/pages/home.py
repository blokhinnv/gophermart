import dash_bootstrap_components as dbc
import dash_loading_spinners as dls
from dash import ctx, dash_table, html, no_update
from dash_extensions.enrich import Input, Output, State, Trigger, callback
from flask_login import current_user
from front.const import ACCRUE, WITHDRAW
from front.app import app

layout = html.Div(
    children=dls.Hash(
        html.Div(
            [
                dbc.Container(
                    [
                        dbc.Row(
                            dbc.Col(
                                [
                                    html.Div(
                                        [
                                            html.Img(
                                                src="/assets/gopher.png",
                                                style={"max-width": 30, "margin-right": 10},
                                            ),
                                            html.H2(
                                                "User",
                                                style={"margin-right": 10},
                                                id="home_username",
                                            ),
                                            html.Div(
                                                [
                                                    html.P(
                                                        [
                                                            html.Span(
                                                                "",
                                                                className="font-weight-bold",
                                                                id="points",
                                                            ),
                                                            "₽",
                                                        ],
                                                        className="m-0",
                                                    )
                                                ],
                                                className="bg-primary text-white px-2 py-2 rounded",
                                                id="points_container",
                                            ),
                                        ],
                                        className="d-flex align-items-center",
                                    )
                                ]
                            ),
                            className="mb-3"
                        ),
                        dbc.Row(
                            [
                                dbc.Col(
                                    [
                                        dbc.Input(
                                            type="input",
                                            id="order_id",
                                            placeholder="Введите номер заказа",
                                            className="mb-3",
                                        ),
                                        dbc.RadioItems(
                                            options=[
                                                {
                                                    "label": "Загрузить заказ",
                                                    "value": ACCRUE,
                                                },
                                                {
                                                    "label": "Списать баллы",
                                                    "value": WITHDRAW,
                                                },
                                            ],
                                            value=ACCRUE,
                                            id="order_action",
                                            className="mb-3",
                                        ),
                                        dbc.InputGroup(
                                            [
                                                dbc.InputGroupText("₽"),
                                                dbc.Input(
                                                    placeholder="Amount",
                                                    type="number",
                                                    id="withdraw_sum",
                                                ),
                                                # dbc.InputGroupText(""),
                                            ],
                                            className="mb-3 d-none",
                                            id="withdraw_sum_group",
                                        ),
                                        html.Div(
                                            id="order_alert",
                                            className="mb-3",
                                        ),
                                        dbc.Button(
                                            "OK",
                                            color="primary",
                                            className="me-1",
                                            id="order_btn",
                                        ),
                                    ],
                                    width=3,
                                ),
                                dbc.Col(
                                    [
                                        dash_table.DataTable(id="orders"),
                                        dbc.Button(
                                            "Обновить",
                                            color="primary",
                                            className="me-1 mt-2",
                                            id="update_orders",
                                            style={"margin-top": 10, "float": "right"}
                                        ),
                                    ]
                                ),
                            ]
                        ),
                    ]
                ),
            ],
            className="mt-3",
        ),
        color="#435278",
        speed_multiplier=2,
        size=100,
        debounce=200,
        fullscreen=True,
    )
)


@callback(
    Output("order_alert", "children"),
    Output("orders", "data"),
    Output("points", "children"),
    Output("points_container", "className"),
    Trigger("order_btn", "n_clicks"),
    State("order_id", "value"),
    State("order_action", "value"),
    State("withdraw_sum", "value"),
)
def process(order_id, order_action, withdraw_sum):
    def gen_outputs(alert_msg, orders, points):
        return (
            dbc.Alert(
                alert_msg,
                color="danger",
                duration=3000,
            )
            if alert_msg
            else no_update,
            orders or no_update,
            int(points) or no_update,
            "bg-primary text-white px-2 py-2 rounded" if points else "d-none",
        )

    # получаем баланс
    balance_res = app.api.balance(current_user.sess)
    if not balance_res["success"]:
        return gen_outputs(balance_res["msg"], None, None)
    current_balance = balance_res["response"]["current"]

    # получаем список заказов
    orders_res = app.api.get_orders(current_user.sess)
    if not orders_res["success"]:
        return gen_outputs(orders_res["msg"], None, current_balance)
    current_orders = orders_res["response"]

    # если это запуск при инициализации, то заканчиваем
    if ctx.triggered_id is None or not order_id:
        return gen_outputs(None, current_orders, current_balance)

    if order_action == ACCRUE:
        # добавляем заказ
        post_order_res = app.api.post_order(order_id, current_user.sess)
        if not post_order_res["success"]:
            return gen_outputs(post_order_res["msg"], current_orders, current_balance)

        # обновляем список заказов
        updated_orders_res = app.api.get_orders(current_user.sess)
        if not updated_orders_res["success"]:
            return gen_outputs(
                updated_orders_res["msg"], current_orders, current_balance
            )
        return gen_outputs(None, updated_orders_res["response"], current_balance)
    else:
        # снимаем баллы
        withdraw_res = app.api.withdraw(order_id, withdraw_sum, current_user.sess)
        if not withdraw_res["success"]:
            return gen_outputs(withdraw_res["msg"], current_orders, current_balance)
        # обновляем баланс
        updated_balance_res = app.api.balance(current_user.sess)
        if not updated_balance_res["success"]:
            return gen_outputs(updated_balance_res["msg"], current_orders, None)
        return gen_outputs(
            None, current_orders, updated_balance_res["response"]["current"]
        )


@callback(Output("withdraw_sum_group", "className"), Input("order_action", "value"))
def change_visibility(order_action):
    if order_action == ACCRUE:
        return "mb-3 d-none"
    return "mb-3 d-flex"


@callback(Output("update_orders", "className"), Input("orders", "data"))
def change_visibility_btn(orders):
    if orders:
        return "me-1 btn btn-primary d-block"
    return "me-1 btn btn-primary d-none"


@callback(
    Output("orders", "data"),
    Output("points", "children"),
    Output("points_container", "className"),
    Trigger("update_orders", "n_clicks"),
    prevent_initial_call=True,
    priority=1,
)
def update():
    orders_res = app.api.get_orders(current_user.sess)
    if not orders_res["success"]:
        return no_update, no_update, no_update

    balance_res = app.api.balance(current_user.sess)
    if not balance_res["success"]:
        return no_update, no_update, no_update
    current_balance = int(balance_res["response"]["current"])

    return (
        orders_res["response"],
        current_balance,
        "bg-primary text-white px-2 py-2 rounded",
    )


@callback(
    Output("home_username", "children"),
    Trigger("update_orders", "n_clicks"),
)
def update_username():
    return current_user.username
