from coolinary.common.v1 import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListOffersRequest(_message.Message):
    __slots__ = ("postcode", "country_code", "page")
    POSTCODE_FIELD_NUMBER: _ClassVar[int]
    COUNTRY_CODE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    postcode: str
    country_code: str
    page: _common_pb2.PageRequest
    def __init__(self, postcode: _Optional[str] = ..., country_code: _Optional[str] = ..., page: _Optional[_Union[_common_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class ListOffersResponse(_message.Message):
    __slots__ = ("offers", "page")
    OFFERS_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    offers: _containers.RepeatedCompositeFieldContainer[Offer]
    page: _common_pb2.PageResponse
    def __init__(self, offers: _Optional[_Iterable[_Union[Offer, _Mapping]]] = ..., page: _Optional[_Union[_common_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class Offer(_message.Message):
    __slots__ = ("id", "shop_id", "product_name", "price", "conditions", "valid_from", "valid_to", "source_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    SHOP_ID_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    CONDITIONS_FIELD_NUMBER: _ClassVar[int]
    VALID_FROM_FIELD_NUMBER: _ClassVar[int]
    VALID_TO_FIELD_NUMBER: _ClassVar[int]
    SOURCE_URL_FIELD_NUMBER: _ClassVar[int]
    id: str
    shop_id: str
    product_name: str
    price: _common_pb2.Money
    conditions: str
    valid_from: str
    valid_to: str
    source_url: str
    def __init__(self, id: _Optional[str] = ..., shop_id: _Optional[str] = ..., product_name: _Optional[str] = ..., price: _Optional[_Union[_common_pb2.Money, _Mapping]] = ..., conditions: _Optional[str] = ..., valid_from: _Optional[str] = ..., valid_to: _Optional[str] = ..., source_url: _Optional[str] = ...) -> None: ...
