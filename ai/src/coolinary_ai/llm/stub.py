"""Rule-based stand-in for an LLM, so the whole pipeline (gRPC, cache, tracing, metrics)
runs before a provider is chosen. Handles simple English, Polish and German lists."""

import re

from .base import ExtractResult, Item, Unit

NUMBER_WORDS = {
    "a": 1,
    "an": 1,
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "dozen": 12,
    "jeden": 1,
    "jedna": 1,
    "jedno": 1,
    "dwa": 2,
    "dwie": 2,
    "trzy": 3,
    "cztery": 4,
    "pięć": 5,
    "sześć": 6,
    "ein": 1,
    "eine": 1,
    "zwei": 2,
    "drei": 3,
    "vier": 4,
    "fünf": 5,
    "sechs": 6,
}
FILLER = {
    "bought",
    "got",
    "some",
    "and",
    "i",
    "the",
    "kupiłem",
    "kupiłam",
    "trochę",
    "oraz",
    "gekauft",
    "ich",
    "habe",
    "etwas",
    "und",
}
SEPARATORS = re.compile(r"[,;\n]+|\s+(?:and|i|oraz|und)\s+", re.IGNORECASE)
AMOUNT = re.compile(r"^(\d+(?:[.,]\d+)?)(g|kg|ml|l)?$", re.IGNORECASE)
UNITS = {
    "g": (Unit.GRAM, 1),
    "kg": (Unit.GRAM, 1000),
    "ml": (Unit.MILLILITRE, 1),
    "l": (Unit.MILLILITRE, 1000),
}


class StubProvider:
    name = "stub"

    async def extract_items(self, text: str, language: str) -> ExtractResult:
        items = [item for chunk in SEPARATORS.split(text) if (item := _parse(chunk))]
        return ExtractResult(
            items=items, model="stub-rules-v1", usage={"input": len(text.split()), "output": len(items)}
        )


def _parse(chunk: str) -> Item | None:
    words = [w for w in chunk.strip().split() if w.lower() not in FILLER]
    if not words:
        return None
    quantity, unit = None, Unit.UNSPECIFIED
    first = words[0].lower()
    if m := AMOUNT.match(first):
        quantity = float(m.group(1).replace(",", "."))
        words = words[1:]
        # "500g flour" or "500 g flour"
        symbol = (m.group(2) or (words[0] if words else "")).lower()
        if symbol in UNITS:
            unit, factor = UNITS[symbol]
            quantity *= factor
            if not m.group(2):
                words = words[1:]
        else:
            unit = Unit.PIECE
    elif first in NUMBER_WORDS:
        quantity, unit = float(NUMBER_WORDS[first]), Unit.PIECE
        words = words[1:]
    name = " ".join(words).strip(" .")
    if not name:
        return None
    return Item(name=name, quantity=quantity, unit=unit, confidence=0.4)
