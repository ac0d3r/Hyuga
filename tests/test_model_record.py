import unittest
from hyuga.models.record import *
from hyuga.lib.option import CONFIG


class TestModelsRecord(unittest.TestCase):
    """Test hyuga.models.record"""

    def setUp(self):
        """create test data
        """
        for i in range(3):
            uid = "test_%s" % i
            name = "test%s.%s" % (i, CONFIG.DOMAIN)
            DnsRecord.objects.create(
                uidentify=uid,
                name=name,
                remote_addr="127.0.0.1"
            )

        for i in range(3):
            uid = "test_%s" % i
            name = "http://test%s.%s" % (i, CONFIG.DOMAIN)
            HttpRecord.objects.create(
                uidentify=uid,
                name=name,
                method="GET",
                date=None,
                user_agent=None,
                content_type=None,
                remote_addr="127.0.0.1"
            )

    def tearDown(self):
        """clear test data
        """
        for dns in DnsRecord.objects.all():
            if dns:
                dns.delete()
        for http in HttpRecord.objects.all():
            if http:
                http.delete()

    def test_dns_record(self):
        uid = "test_%s" % 1
        name = "test%s.%s" % (1, CONFIG.DOMAIN)
        dns = DnsRecord.objects.filter(uidentify=uid)
        self.assertEqual(len(dns), 1)
        self.assertEqual(dns[0].name, name)

    def test_dns_record_errors(self):
        d = DnsRecord(
            uidentify="test_test",
            name=CONFIG.DOMAIN,
            remote_addr="127.0.0.1"
        )
        self.assertFalse(d.is_valid())
        self.assertIn(("name", 'Not DnsRecord name'), d.errors)

        d1 = DnsRecord(
            uidentify="test_test",
            name="test%s.%s" % (10, CONFIG.DOMAIN)
        )
        self.assertFalse(d1.is_valid())
        self.assertIn(('remote_addr', 'required'), d1.errors)

    def test_http_record(self):
        uid = "test_%s" % 1
        name = "http://test%s.%s" % (1, CONFIG.DOMAIN)
        http = HttpRecord.objects.filter(uidentify=uid)
        self.assertEqual(len(http), 1)
        self.assertEqual(http[0].name, name)

    def test_http_record_errors(self):
        # 测试创建不符合要求的数据
        h = HttpRecord(
            uidentify="test_test",
            name="http://hyuga.io",
            method="GET",
            remote_addr="127.0.0.1"
        )
        self.assertFalse(h.is_valid())
        self.assertIn(("name", 'Not HttpRecord name'), h.errors)

        h = HttpRecord(
            uidentify="test_test",
            name="http://test.hyuga.io/",
            remote_addr="127.0.0.1"
        )
        self.assertFalse(h.is_valid())
        self.assertIn(("method", 'required'), h.errors)

    def test_validator(self):
        field_name = "name"
        hdr = hyuga_domain_validator(field_name, CONFIG.DOMAIN)
        self.assertEqual(hdr, ((field_name, 'Not DnsRecord name'),))

        hur = hyuga_url_validator(field_name, "http://%s/abcd" % CONFIG.DOMAIN)
        self.assertEqual(hur, ((field_name, 'Not HttpRecord name'),))

        hdr = hyuga_domain_validator(field_name, "abc.%s" % CONFIG.DOMAIN)
        self.assertEqual(hdr, None)

        hur = hyuga_url_validator(
            field_name, "http://abc.%s/abcd" % CONFIG.DOMAIN)
        self.assertEqual(hur, None)


if __name__ == '__main__':
    unittest.main()
