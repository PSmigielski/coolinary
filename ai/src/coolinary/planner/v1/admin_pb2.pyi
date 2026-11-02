from coolinary.common.v1 import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListAdminChangesRequest(_message.Message):
    __slots__ = ("target", "page")
    TARGET_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    target: str
    page: _common_pb2.PageRequest
    def __init__(self, target: _Optional[str] = ..., page: _Optional[_Union[_common_pb2.PageRequest, _Mapping]] = ...) -> None: ...

class AdminChange(_message.Message):
    __slots__ = ("admin_user_id", "rpc", "target", "summary_json", "created_at")
    ADMIN_USER_ID_FIELD_NUMBER: _ClassVar[int]
    RPC_FIELD_NUMBER: _ClassVar[int]
    TARGET_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_JSON_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    admin_user_id: str
    rpc: str
    target: str
    summary_json: str
    created_at: str
    def __init__(self, admin_user_id: _Optional[str] = ..., rpc: _Optional[str] = ..., target: _Optional[str] = ..., summary_json: _Optional[str] = ..., created_at: _Optional[str] = ...) -> None: ...

class ListAdminChangesResponse(_message.Message):
    __slots__ = ("changes", "page")
    CHANGES_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    changes: _containers.RepeatedCompositeFieldContainer[AdminChange]
    page: _common_pb2.PageResponse
    def __init__(self, changes: _Optional[_Iterable[_Union[AdminChange, _Mapping]]] = ..., page: _Optional[_Union[_common_pb2.PageResponse, _Mapping]] = ...) -> None: ...

class GetQueueCountsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetQueueCountsResponse(_message.Message):
    __slots__ = ("counts", "computed_at")
    class CountsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    COUNTS_FIELD_NUMBER: _ClassVar[int]
    COMPUTED_AT_FIELD_NUMBER: _ClassVar[int]
    counts: _containers.ScalarMap[str, int]
    computed_at: str
    def __init__(self, counts: _Optional[_Mapping[str, int]] = ..., computed_at: _Optional[str] = ...) -> None: ...
