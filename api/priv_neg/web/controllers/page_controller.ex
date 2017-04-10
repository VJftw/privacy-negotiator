defmodule PrivNeg.PageController do
  use PrivNeg.Web, :controller

  def index(conn, _params) do
    conn
    |> put_status(200)
    |> render("index.json", %{})
  end
end
