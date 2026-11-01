from coolinary.common.v1 import common_pb2 as _common_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ResolveSource(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RESOLVE_SOURCE_UNSPECIFIED: _ClassVar[ResolveSource]
    RESOLVE_SOURCE_HOUSEHOLD: _ClassVar[ResolveSource]
    RESOLVE_SOURCE_SHARED_CATALOGUE: _ClassVar[ResolveSource]
    RESOLVE_SOURCE_OPEN_FOOD_FACTS: _ClassVar[ResolveSource]
    RESOLVE_SOURCE_GS1: _ClassVar[ResolveSource]
RESOLVE_SOURCE_UNSPECIFIED: ResolveSource
RESOLVE_SOURCE_HOUSEHOLD: ResolveSource
RESOLVE_SOURCE_SHARED_CATALOGUE: ResolveSource
RESOLVE_SOURCE_OPEN_FOOD_FACTS: ResolveSource
RESOLVE_SOURCE_GS1: ResolveSource

class ResolveBarcodeRequest(_message.Message):
    __slots__ = ("ean", "household_id")
    EAN_FIELD_NUMBER: _ClassVar[int]
    HOUSEHOLD_ID_FIELD_NUMBER: _ClassVar[int]
    ean: str
    household_id: str
    def __init__(self, ean: _Optional[str] = ..., household_id: _Optional[str] = ...) -> None: ...

class ResolveBarcodeResponse(_message.Message):
    __slots__ = ("found", "product", "source")
    FOUND_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    found: bool
    product: Product
    source: ResolveSource
    def __init__(self, found: _Optional[bool] = ..., product: _Optional[_Union[Product, _Mapping]] = ..., source: _Optional[_Union[ResolveSource, str]] = ...) -> None: ...

class Product(_message.Message):
    __slots__ = ("ean", "name", "brand", "pack_size")
    EAN_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    BRAND_FIELD_NUMBER: _ClassVar[int]
    PACK_SIZE_FIELD_NUMBER: _ClassVar[int]
    ean: str
    name: str
    brand: str
    pack_size: _common_pb2.Quantity
    def __init__(self, ean: _Optional[str] = ..., name: _Optional[str] = ..., brand: _Optional[str] = ..., pack_size: _Optional[_Union[_common_pb2.Quantity, _Mapping]] = ...) -> None: ...
