#!/usr/bin/python3
# -*- coding: utf-8 -*-

import socket
import mysql.connector
import sys

mysqldb = mysql.connector.connect(
        host="avaya_db",
        user="cdruser",
        passwd="averysecurepasswordforever",
        db="AvayaCdr"
)


sock = socket.socket()
sock.bind(('', 60010))
sock.listen(1)
conn, addr = sock.accept()

mysqldbcursor = mysqldb.cursor()

while True:

    try:
        data = conn.recv(1024)
        if not data:
            break
        cdr = data.decode("utf-8").rstrip().split(",")

        list = [cdr[0], cdr[1], cdr[2], cdr[3], cdr[4], cdr[5], cdr[6], cdr[8], cdr[9], cdr[15], cdr[16], cdr[27], cdr[28], cdr[29], cdr[30], cdr[31] , cdr[33], cdr[34]]
        print("[INFO]: " + ' '.join(map(str, list)))

        sql_insert = "INSERT INTO cdr (CallStart,ConnectedTime,RingTime,Caller,Direction,CalledNumber,DialledNumber,IsInternal,CallID,HoldTime,ParkTime,ExternalTargetingCause,ExternalTargeterId,ExternalTargetedNumber,CallerServerIP,UniqueCallIDCallerExtension,UniqueCallIDCalledParty,SMDRRecordingTime) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"

        try:
            mysqldbcursor.execute(sql_insert, list)
            mysqldb.commit()
        except:
            print("[ERROR]: Could not add record " + ' '.join(map(str, list)))
    except KeyboardInterrupt:
        sys.exit()

mysqldbcursor.close()
mysqldb.close()

conn.close()