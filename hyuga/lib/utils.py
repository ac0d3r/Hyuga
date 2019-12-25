import random
import string

import shortuuid


def rand_str(count):
    return ''.join([random.choice(string.ascii_letters + string.digits) for i in range(count)])


def get_shortuuid(length=6) -> str:
    return shortuuid.ShortUUID().random(length).lower()


def records_to_list(records, filter_=None):
    """
    :type records: Iterable
    :type filter_: str
    :return: List
    """
    _data = []
    for record in records:
        if not record:
            continue
        if not (filter_ is None):
            if not (isinstance(record.name, str) and filter_ in record.name):
                continue
        try:
            _data.append(record.to_dict())
        except:
            pass
    return _data
