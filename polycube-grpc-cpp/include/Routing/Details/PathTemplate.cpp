#include "PathTemplate.h"

namespace Routing {
namespace Details {

PathTemplate::PathTemplate(std::string const &path, std::string const &http_verb)
: _parser(path), _templ(path), _http_verb(http_verb) {
    _regex = RegexConverter().toRegex(path);
}

}
}