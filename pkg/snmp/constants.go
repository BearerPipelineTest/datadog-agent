package snmp

// DeviceReachableGetNextOid is used in getNext call to check if the device is reachable
// GETNEXT 1.0 should be able to fetch the first available SNMP OID.
// There is no need to handle top node other than iso(1) since it the only valid SNMP tree starting point.
// Other top nodes like ccitt(0) and joint(2) do not pertain to SNMP.
// Source: https://docstore.mik.ua/orelly/networking_2ndEd/snmp/ch02_03.htm
const DeviceReachableGetNextOid = "1.0"
