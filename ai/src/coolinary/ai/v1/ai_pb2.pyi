from coolinary.common.v1 import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ExtractItemsFromTextRequest(_message.Message):
    __slots__ = ("text", "language", "household_id")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_FIELD_NUMBER: _ClassVar[int]
    HOUSEHOLD_ID_FIELD_NUMBER: _ClassVar[int]
    text: str
    language: str
    household_id: str
    def __init__(self, text: _Optional[str] = ..., language: _Optional[str] = ..., household_id: _Optional[str] = ...) -> None: ...

class ExtractItemsFromTextResponse(_message.Message):
    __slots__ = ("items", "model", "cached")
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    CACHED_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[ExtractedItem]
    model: str
    cached: bool
    def __init__(self, items: _Optional[_Iterable[_Union[ExtractedItem, _Mapping]]] = ..., model: _Optional[str] = ..., cached: _Optional[bool] = ...) -> None: ...

class ExtractedItem(_message.Message):
    __slots__ = ("name", "quantity", "confidence")
    NAME_FIELD_NUMBER: _ClassVar[int]
    QUANTITY_FIELD_NUMBER: _ClassVar[int]
    CONFIDENCE_FIELD_NUMBER: _ClassVar[int]
    name: str
    quantity: _common_pb2.Quantity
    confidence: float
    def __init__(self, name: _Optional[str] = ..., quantity: _Optional[_Union[_common_pb2.Quantity, _Mapping]] = ..., confidence: _Optional[float] = ...) -> None: ...
