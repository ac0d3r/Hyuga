import falcon
import redis

from hyuga.api.common import FIELDS, BaseResource, BaseValidate
from hyuga.core import errors
from hyuga.core.auth import authenticated
from hyuga.core.log import _api_logger
from hyuga.lib import utils
from hyuga.lib.option import CONFIG
from hyuga.models.record import DnsRecord, HttpRecord
from hyuga.models.user import User

RECORD = {
    "dns": DnsRecord,
    "http": HttpRecord
}


class UsersSelfRecords(BaseResource):
    """
    Handle for endpoint: /v1/users/self/records
    """
    Schema = {
        "type": FIELDS["recordType"]
    }

    @authenticated()
    @falcon.before(BaseValidate(Schema, is_params=True).validate)
    def on_get(self, req, resp):
        Record = RECORD[req.params["type"]]
        try:
            records = Record.objects.filter(
                uidentify=self.current_user.identify).limit(CONFIG.RECORDS_MAX_NUM)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersSelfRecords on_get ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

        self.on_success(resp, utils.records_to_list(records))

    @authenticated()
    @falcon.before(BaseValidate(Schema).validate)
    def on_delete(self, req, resp):
        Record = RECORD[req.context["data"]["type"]]

        try:
            records = Record.objects.filter(
                uidentify=self.current_user.identify)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersSelfRecords on_delete ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

        if records:
            UsersSelfRecords.delete_records(records)
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
        _filter = None
        _token = req.params["token"]
        if "filter" in req.params.keys():
            _filter = req.params["filter"]

        Record = RECORD[req.params["type"]]
        try:
            user = User.objects.filter(token=_token).first()
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("Records on_get ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

        if not user:
            raise errors.UnauthorizedError("Token Does Not Exist: %s" % _token)

        try:
            records = Record.objects.filter(
                uidentify=user.identify).limit(CONFIG.RECORDS_MAX_NUM)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("Records on_get ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

        resp_data = utils.records_to_list(records, _filter)
        self.on_success(resp, resp_data)
