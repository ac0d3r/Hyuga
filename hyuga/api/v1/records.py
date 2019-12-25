import falcon

from hyuga.api.common import FIELDS, BaseResource, BaseValidate
from hyuga.core.auth import authenticated
from hyuga.core.errors import (InvalidParameterError, UnauthorizedError,
                               UserNotExistsError)
from hyuga.core.log import _api_logger
from hyuga.lib.option import CONFIG
from hyuga.lib.utils import records_to_list
from hyuga.models.record import DnsRecord, HttpRecord
from hyuga.models.user import User


class UsersSelfRecords(BaseResource):
    """
    Handle for endpoint: /v1/users/self/records
    """
    Schema = {
        "type": FIELDS["recordType"]
    }

    @falcon.before(BaseValidate(Schema, is_params=True).validate)
    @authenticated
    def on_get(self, req, resp):
        if self.current_user is None:
            raise UnauthorizedError()
        if not req.params:
            raise InvalidParameterError()

        Record = None
        _type = req.params["type"]
        if _type == "dns":
            Record = DnsRecord
        elif _type == "http":
            Record = HttpRecord

        records = Record.objects.filter(
            uidentify=self.current_user.identify).limit(CONFIG.RECORDS_MAX_NUM)
        resp_data = records_to_list(records)
        self.on_success(resp, resp_data)

    @falcon.before(BaseValidate(Schema).validate)
    @authenticated
    def on_delete(self, req, resp):
        if self.current_user is None:
            raise UnauthorizedError()
        req_data = req.context["data"]
        if not req_data:
            raise InvalidParameterError()

        Record = None
        _type = req_data["type"]
        if _type == "dns":
            Record = DnsRecord
        elif _type == "http":
            Record = HttpRecord

        records = Record.objects.filter(uidentify=self.current_user.identify)
        try:
            UsersSelfRecords.delete_records(records)
        except:
            pass
        self.on_success(resp)

    @staticmethod
    def delete_records(records):
        """
        :type records: Iterable
        """
        for record in records:
            record.delete()


class Records(BaseResource):
    """
    Handle for endpoint: /v1/records
    """
    getSchema = {
        "type": FIELDS["recordType"],
        "token": FIELDS["userToken"],
        "filter": FIELDS["recordsFilter"]
    }

    @falcon.before(BaseValidate(getSchema, is_params=True).validate)
    def on_get(self, req, resp):
        if not req.params:
            raise InvalidParameterError()

        _filter = None
        _type = req.params["type"]
        _token = req.params["token"]
        if "filter" in req.params.keys():
            _filter = req.params["filter"]

        Record = None
        _type = req.params["type"]
        if _type == "dns":
            Record = DnsRecord
        elif _type == "http":
            Record = HttpRecord

        try:
            user = User.get(User.token == _token)
            records = Record.objects.filter(
                uidentify=user.identify).limit(CONFIG.RECORDS_MAX_NUM)
            resp_data = records_to_list(records, _filter)
            self.on_success(resp, resp_data)
        except User.DoesNotExist:
            raise UnauthorizedError(f"Token Does Not Exist: {_token}")
