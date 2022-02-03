# Simple C++ URI routing library

This library is used internally to the client-side polycube-grpc-cpp library in order to be able to route user requests. The version used actually has some changes compared to the basic version. Here you find the [original](https://github.com/qoollo/simple-cpp-router).


## Requirements
C++11 compiler with regex support

## Example usage

```cpp

    Router router;
    
    router.registerPath("/simple_path");
    router.registerPath("/one/two");
    router.registerPath("/path/with/:dynamic_var");
    router.registerPath("/path/to/:one_var/and/:another_var");
    
    PathMatch match = router.matchPath("/path/to/dynvar_1/and/dynvar_2");
    
    assert(match.pathTemplate() == "/path/to/:one_var/and/:another_var");
    assert(std::string("dynvar_1") == match["one_var"]);
    assert(std::string("dynvar_2") == match["another_var"]);
    
```
