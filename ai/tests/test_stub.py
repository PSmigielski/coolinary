from coolinary_ai.llm import Unit
from coolinary_ai.llm.stub import StubProvider


async def test_extracts_english_list():
    res = await StubProvider().extract_items("bought eggs, two peppers, some cheese and 500 g flour", "en")
    got = {(i.name, i.quantity, i.unit) for i in res.items}
    assert got == {
        ("eggs", None, Unit.UNSPECIFIED),
        ("peppers", 2.0, Unit.PIECE),
        ("cheese", None, Unit.UNSPECIFIED),
        ("flour", 500.0, Unit.GRAM),
    }


async def test_extracts_polish_and_litres():
    res = await StubProvider().extract_items("kupiłem dwie papryki i 1 l mleka", "pl")
    assert [(i.name, i.quantity, i.unit) for i in res.items] == [
        ("papryki", 2.0, Unit.PIECE),
        ("mleka", 1000.0, Unit.MILLILITRE),
    ]
