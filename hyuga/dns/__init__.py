# coding: utf-8
"""Simple DNS Server
"""
import copy
import os
import re
import sys
import tempfile

from dnslib import QTYPE, RCODE, RR, TXT
from dnslib.server import BaseResolver, DNSServer

from hyuga.core.log import _dns_logger
from hyuga.lib.option import CONFIG
from hyuga.models.record import DnsRecord

ZONE = '''
*.{dnsdomain}.       IN      NS      {ns1domain}.
*.{dnsdomain}.       IN      NS      {ns2domain}.
*.{dnsdomain}.       IN      A       {serverip}
{dnsdomain}.       IN      A       {serverip}
'''.format(
    dnsdomain=CONFIG.DOMAIN,
    ns1domain=CONFIG.NS1_DOMAIN,
    ns2domain=CONFIG.NS2_DOMAIN,
    serverip=CONFIG.SERVER_IP
)


class RedisLogger():
    def log_data(self, dnsobj):
        pass

    def log_error(self, handler, e):
        pass

    def log_pass(self, *args):
        pass

    def log_prefix(self, handler):
        pass

    def log_recv(self, handler, data):
        pass

    def log_reply(self, handler, reply):
        pass

    def log_request(self, handler, request):
        client_ip = handler.client_address[0]
        domain = request.q.qname.__str__().lower()

        if not domain.endswith(CONFIG.DOMAIN + '.'):
            return
        r_udomain = re.search(r'\.?([^\.]+\.%s)\.' %
                              CONFIG.DOMAIN, domain)
        if not r_udomain:
            return

        udomain = r_udomain.group(1)
        if udomain in CONFIG.DNS_IGNORE_DOMAINS:
            return

        # 写入 redis
        try:
            dns_record = DnsRecord(
                uidentify=udomain.split("."+CONFIG.DOMAIN)[0],
                name=udomain,
                remote_addr=client_ip
            )
            dns_record.save()
            dns_record.expire(CONFIG.RECORDS_EXPIRE)
        except:
            pass

    def log_send(self, handler, data):
        pass

    def log_truncated(self, handler, reply):
        pass


class ZoneResolver(BaseResolver):
    """
        Simple fixed zone file resolver.
    """

    def __init__(self, zone, glob=False):
        """
            Initialise resolver from zone file.
            Stores RRs as a list of (label,type,rr) tuples
            If 'glob' is True use glob match against zone file
        """
        self.zone = [(rr.rname, QTYPE[rr.rtype], rr)
                     for rr in RR.fromZone(zone)]
        self.glob = glob
        self.eq = 'matchGlob' if glob else '__eq__'

    def resolve(self, request, handler):
        """
            Respond to DNS request - parameters are request packet & handler.
            Method is expected to return DNS response
        """
        reply = request.reply()
        qname = request.q.qname
        qtype = QTYPE[request.q.qtype]
        if qtype == 'TXT':
            txtpath = os.path.join(tempfile.gettempdir(), str(qname).lower())
            if os.path.isfile(txtpath):
                reply.add_answer(
                    RR(qname, QTYPE.TXT, rdata=TXT(open(txtpath).read().strip())))
        for name, rtype, rr in self.zone:
            # Check if label & type match
            if getattr(qname,
                       self.eq)(name) and (qtype == rtype or qtype == 'ANY'
                                           or rtype == 'CNAME'):
                # If we have a glob match fix reply label
                if self.glob:
                    a = copy.copy(rr)
                    a.rname = qname
                    reply.add_answer(a)
                else:
                    reply.add_answer(rr)
                # Check for A/AAAA records associated with reply and
                # add in additional section
                if rtype in ['CNAME', 'NS', 'MX', 'PTR']:
                    for a_name, a_rtype, a_rr in self.zone:
                        if a_name == rr.rdata.label and a_rtype in [
                                'A', 'AAAA'
                        ]:
                            reply.add_ar(a_rr)
        if not reply.rr:
            reply.header.rcode = RCODE.NXDOMAIN
        return reply


def main():
    try:
        resolver = ZoneResolver(ZONE, True)
        logger = RedisLogger()
        _dns_logger.info(ZONE)
        _dns_logger.info("Starting Dns Server (*:53) [UDP]")

        dns_server = DNSServer(resolver, port=53, address='', logger=logger)
        dns_server.start()
    except KeyboardInterrupt:
        _dns_logger.info("Shutting Down Dns Server (*:53) [UDP]")
        dns_server.stop()


if __name__ == '__main__':
    main()
