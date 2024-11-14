# Avaya SMDR to MySQL
A small utility which listens for SMDR call data from one or more Avaya PBX systems and writes the call data to a MySQL database. 
Includes docker config to run as containers:
1. SMDR Receiver to listen for inbound SMDR messages on a specified port and write to a specified database
2. MySQL DB to hold call details in a more useful format than the archaic SMDR standard
