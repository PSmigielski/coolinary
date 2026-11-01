from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Role(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ROLE_UNSPECIFIED: _ClassVar[Role]
    ROLE_OWNER: _ClassVar[Role]
    ROLE_MEMBER: _ClassVar[Role]
    ROLE_PROFILE: _ClassVar[Role]

class Plan(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PLAN_UNSPECIFIED: _ClassVar[Plan]
    PLAN_FREE: _ClassVar[Plan]
    PLAN_PLUS: _ClassVar[Plan]
ROLE_UNSPECIFIED: Role
ROLE_OWNER: Role
ROLE_MEMBER: Role
ROLE_PROFILE: Role
PLAN_UNSPECIFIED: Plan
PLAN_FREE: Plan
PLAN_PLUS: Plan

class GetHouseholdRequest(_message.Message):
    __slots__ = ("household_id",)
    HOUSEHOLD_ID_FIELD_NUMBER: _ClassVar[int]
    household_id: str
    def __init__(self, household_id: _Optional[str] = ...) -> None: ...

class GetHouseholdResponse(_message.Message):
    __slots__ = ("household",)
    HOUSEHOLD_FIELD_NUMBER: _ClassVar[int]
    household: Household
    def __init__(self, household: _Optional[_Union[Household, _Mapping]] = ...) -> None: ...

class Household(_message.Message):
    __slots__ = ("id", "name", "owner_user_id", "members", "plan")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    OWNER_USER_ID_FIELD_NUMBER: _ClassVar[int]
    MEMBERS_FIELD_NUMBER: _ClassVar[int]
    PLAN_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    owner_user_id: str
    members: _containers.RepeatedCompositeFieldContainer[Member]
    plan: Plan
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ..., owner_user_id: _Optional[str] = ..., members: _Optional[_Iterable[_Union[Member, _Mapping]]] = ..., plan: _Optional[_Union[Plan, str]] = ...) -> None: ...

class Member(_message.Message):
    __slots__ = ("id", "display_name", "role", "has_account")
    ID_FIELD_NUMBER: _ClassVar[int]
    DISPLAY_NAME_FIELD_NUMBER: _ClassVar[int]
    ROLE_FIELD_NUMBER: _ClassVar[int]
    HAS_ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    id: str
    display_name: str
    role: Role
    has_account: bool
    def __init__(self, id: _Optional[str] = ..., display_name: _Optional[str] = ..., role: _Optional[_Union[Role, str]] = ..., has_account: _Optional[bool] = ...) -> None: ...
