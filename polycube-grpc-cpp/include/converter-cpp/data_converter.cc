#include "data_converter.h"


namespace converter {

uint16_t bytesToUint16_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 2){
    std::cerr << "Error: data is to small";
    return 0;
  }

  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right     
    return (data[start] << 8) | (data[start+1]);
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

    return (data[start]) | (data[start+1] << 8);
  }
}


int16_t bytesToInt16_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 2){
    std::cerr << "Error: data is to small";
    return 0;
  }

  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right     
    return (data[start] << 8) | (data[start+1]);
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

    return (data[start]) | (data[start+1] << 8);
  }
}


uint32_t bytesToUint32_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 4){
    std::cerr << "Error: data is to small";
    return 0;
  }

  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right     
    return (data[start] << 24) | (data[start+1] << 16) |
      (data[start+2] << 8) | (data[start+3]);
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

    return (data[start]) | (data[start+1] << 8) |
      (data[start+2] << 16) | (data[start+3] << 24);
  }
}



int32_t bytesToInt32_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 4){
    std::cerr << "Error: data is to small";
    return 0;
  }

  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right     
    return (data[start] << 24) | (data[start+1] << 16) |
      (data[start+2] << 8) | (data[start+3]);
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

    return (data[start]) | (data[start+1] << 8) |
      (data[start+2] << 16) | (data[start+3] << 24);
  }
}

  



uint64_t bytesToUint64_t_48(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 6){
    std::cerr << "Error: data is to small";
    return 0;
  }
  std::cout << "800A" << std::endl;
  for(int i=0;i<6;i++)
    printf("%d - %.2x \n", i, data[start+i]);
  std::cout << std::endl;
      
  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right  
    std::cout <<  "bytesToUint64_t_48 isNetByteOrder" << std::endl;   
    uint32_t a = ((uint8_t)(data[start]) << 16 | (uint8_t)(data[start+1]) << 8 | (uint8_t)(data[start+2]));
    uint32_t b = ((uint8_t)(data[start+3]) << 16 | (uint8_t)(data[start+4]) << 8 | (uint8_t)(data[start+5]));
    return int64_t((uint64_t) a << 24 | ((uint32_t) b));
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

  std::cout <<  "bytesToUint64_t_48 NOT isNetByteOrder" << std::endl;   

    uint32_t a = ((uint8_t)(data[start]) | (uint8_t)(data[start+1]) << 8 | (uint8_t)(data[start+2]) << 16);
    uint32_t b = ((uint8_t)(data[start+3]) | (uint8_t)(data[start+4]) << 8 | (uint8_t)(data[start+5]) << 16);
    return int64_t((uint32_t) a | ((uint64_t) b << 24));
  }
}

/*
Why the conversion is divided in two parts?
Possible overflow, detected also by the compiler as a warning when we try to << 40
*/
uint64_t bytesToUint64_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 8){
    std::cerr << "Error: data is to small";
    return 0;
  }
  std::cout << "800A" << std::endl;
  std::cout << data[start] << " - " << data[start+1] << " - "  
      << data[start+2] << " - " << data[start+3] << " - "
      << data[start+4] << " - " << data[start+5] << " - "
      << data[start+6] << " - " << data[start+7] << " - \n  ";
      
  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right 
    std::cout <<  "bytesToUint64_t isNetByteOrder" << std::endl;   
    uint32_t a = ((uint8_t)(data[start]) << 24 | (uint8_t)(data[start+1]) << 16 | (uint8_t)(data[start+2]) << 8) | (uint8_t)(data[start+3]);
    uint32_t b = ((uint8_t)(data[start+4]) << 24 | (uint8_t)(data[start+5]) << 16 | (uint8_t)(data[start+6]) << 8) | (uint8_t)(data[start+7]);
    return int64_t((uint64_t) a << 32 | ((uint32_t) b));
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

   std::cout <<  "bytesToUint64_t NOT isNetByteOrder" << std::endl;   
    uint32_t a = ((uint8_t)(data[start]) | (uint8_t)(data[start+1]) << 8 | (uint8_t)(data[start+2]) << 16) | ((uint8_t)(data[start+3]) << 24);
    uint32_t b = ((uint8_t)(data[start+4]) | (uint8_t)(data[start+5]) << 8 | (uint8_t)(data[start+6]) << 16) | ((uint8_t)(data[start+7]) << 24);
    return int64_t((uint32_t) a | ((uint64_t) b << 32));
  }
}


/*
CONTROL IF IT WORKS: the unsigned version works
*/
int64_t bytesToInt64_t(const std::string& data, int start) {

  if(data.empty()){
    std::cerr << "Error: data is empty";
    return 0;
  }
  if(start > data.size() - 8){
    std::cerr << "Error: data is to small";
    return 0;
  }

  if(isNetByteOrder){  // Big endian, Network byte order, MSB on lower address, the MSB is the first written from left ro right 
    std::cout <<  "bytesToInt64_t isNetByteOrder" << std::endl;   
    int32_t a = ((int8_t)(data[start]) << 24 | (int8_t)(data[start+1]) << 16 | (int8_t)(data[start+2]) << 8) | (int8_t)(data[start+3]);
    int32_t b = ((int8_t)(data[start+4]) << 24 | (int8_t)(data[start+5]) << 16 | (int8_t)(data[start+6]) << 8) | (int8_t)(data[start+7]);
    return int64_t((int64_t) a << 32 | ((int32_t) b));
  } 
  else {  // Little endian, Host byte order, MSB on higher address, the MSB is the last written from left to right

   std::cout <<  "bytesToInt64_t NOT isNetByteOrder" << std::endl;   
    uint32_t a = ((int8_t)(data[start]) | (int8_t)(data[start+1]) << 8 | (int8_t)(data[start+2]) << 16) | ((int8_t)(data[start+3]) << 24);
    uint32_t b = ((int8_t)(data[start+4]) | (int8_t)(data[start+5]) << 8 | (int8_t)(data[start+6]) << 16) | ((int8_t)(data[start+7]) << 24);
    return int64_t((int32_t) a | ((int64_t) b << 32));
  }
}



}