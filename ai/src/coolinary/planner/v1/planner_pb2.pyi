from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Meal(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    MEAL_UNSPECIFIED: _ClassVar[Meal]
    MEAL_BREAKFAST: _ClassVar[Meal]
    MEAL_DINNER: _ClassVar[Meal]
    MEAL_SUPPER: _ClassVar[Meal]
MEAL_UNSPECIFIED: Meal
MEAL_BREAKFAST: Meal
MEAL_DINNER: Meal
MEAL_SUPPER: Meal

class SuggestDishesRequest(_message.Message):
    __slots__ = ("household_id", "meal", "limit")
    HOUSEHOLD_ID_FIELD_NUMBER: _ClassVar[int]
    MEAL_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    household_id: str
    meal: Meal
    limit: int
    def __init__(self, household_id: _Optional[str] = ..., meal: _Optional[_Union[Meal, str]] = ..., limit: _Optional[int] = ...) -> None: ...

class SuggestDishesResponse(_message.Message):
    __slots__ = ("suggestions",)
    SUGGESTIONS_FIELD_NUMBER: _ClassVar[int]
    suggestions: _containers.RepeatedCompositeFieldContainer[DishSuggestion]
    def __init__(self, suggestions: _Optional[_Iterable[_Union[DishSuggestion, _Mapping]]] = ...) -> None: ...

class DishSuggestion(_message.Message):
    __slots__ = ("dish_id", "name", "share_on_hand", "kcal_per_portion")
    DISH_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    SHARE_ON_HAND_FIELD_NUMBER: _ClassVar[int]
    KCAL_PER_PORTION_FIELD_NUMBER: _ClassVar[int]
    dish_id: str
    name: str
    share_on_hand: float
    kcal_per_portion: int
    def __init__(self, dish_id: _Optional[str] = ..., name: _Optional[str] = ..., share_on_hand: _Optional[float] = ..., kcal_per_portion: _Optional[int] = ...) -> None: ...
