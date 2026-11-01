from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SendEmailRequest(_message.Message):
    __slots__ = ("to", "subject", "text_body", "html_body")
    TO_FIELD_NUMBER: _ClassVar[int]
    SUBJECT_FIELD_NUMBER: _ClassVar[int]
    TEXT_BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_BODY_FIELD_NUMBER: _ClassVar[int]
    to: str
    subject: str
    text_body: str
    html_body: str
    def __init__(self, to: _Optional[str] = ..., subject: _Optional[str] = ..., text_body: _Optional[str] = ..., html_body: _Optional[str] = ...) -> None: ...

class SendEmailResponse(_message.Message):
    __slots__ = ("message_id",)
    MESSAGE_ID_FIELD_NUMBER: _ClassVar[int]
    message_id: str
    def __init__(self, message_id: _Optional[str] = ...) -> None: ...
