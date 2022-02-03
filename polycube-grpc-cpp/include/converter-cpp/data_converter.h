#include <string>
#include <iostream>

namespace converter {

static bool isNetByteOrder = true; // with gRPC we receive data in the same format as the server retrieves them
/*
Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right

Address here means the INDEX in the following representation

Example: hex number 0x0123


      Little endian             Big endian
       +----+----+              +----+----+
       |0x23|0x01|              |0x01|0x23|
       +----+----+              +----+----+
 byte:    0    1                   0    1



*/



/*
Method used to convert a string used to store data as bytes in a uin16_t
*/
uint16_t bytesToUint16_t(const std::string& data, int start);

int16_t bytesToInt16_t(const std::string& data, int start);


/*
* Method used to convert a string used to store data as bytes in a uint32_t
* the start parameter is used ad a starting point, i.e., start = 0, the method 
* uses data[0], data[1], data[2], data[3], with the correct bitshift (24,16,8,0) or (0,8,16,24)
* depending on the NBO or HBO
*/
uint32_t bytesToUint32_t(const std::string& data, int start);

int32_t bytesToInt32_t(const std::string& data, int start);


/*
* These methods are used to convert bytes to uint64 truncated to 48. Used only for the MAC address
*/
uint64_t bytesToUint64_t_48(const std::string& data, int start);

int64_t bytesToInt64_t_48(const std::string& data, int start);






/*
Method used to convert a string used to store data as bytes in a uint64_t
*/
uint64_t bytesToUint64_t(const std::string& data, int start);

int64_t bytesToInt64_t(const std::string& data, int start);




}